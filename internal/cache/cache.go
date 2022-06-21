package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	application "github.com/debugger84/modulus-application"
	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure/v2"
	"log"
	"strconv"
	"time"
)

type Cache[KeyT comparable, ValueT any] struct {
	cache        *bigcache.BigCache
	cacheEnabled bool
	logger       application.Logger
}

type Config struct {
	MaxCacheSizeInMb int
	CacheEnabled     bool
	// time after which entry can be evicted
	LifeTime time.Duration
}

func NewCache[KeyT comparable, ValueT any](config *Config, logger application.Logger) *Cache[KeyT, ValueT] {
	if config.CacheEnabled {
		if config.LifeTime == 0 {
			config.LifeTime = 10 * time.Minute
		}
		hardMaxCacheSize := config.MaxCacheSizeInMb
		if hardMaxCacheSize == 0 {
			hardMaxCacheSize = 1
		}
		cfg := bigcache.Config{
			// number of shards (must be a power of 2)
			Shards: 256,
			// time after which entry can be evicted
			LifeWindow: config.LifeTime,
			// Interval between removing expired entries (clean up).
			// If set to <= 0 then no action is performed.
			// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
			CleanWindow: 5 * time.Second,
			// rps * lifeWindow, used only in initial memory allocation
			MaxEntriesInWindow: 1000,
			// max entry size in bytes, used only in initial memory allocation
			MaxEntrySize: 100,

			StatsEnabled: false,

			// prints information about additional memory allocation
			Verbose: false,
			Hasher:  nil,
			// cache will not allocate more memory than this limit, value in MB
			// if value is reached then the oldest entries can be overridden for the new ones
			// 0 value means no size limit
			HardMaxCacheSize: hardMaxCacheSize,
			// callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			OnRemove:             nil,
			OnRemoveWithMetadata: nil,
			// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A constant representing the reason will be passed through.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			// Ignored if OnRemove is specified.
			OnRemoveWithReason: nil,
			Logger:             nil,
		}

		cache, initErr := bigcache.NewBigCache(cfg)
		if initErr != nil {
			log.Fatal(initErr)
		}

		return &Cache[KeyT, ValueT]{cache: cache, cacheEnabled: true, logger: logger}
	}
	return &Cache[KeyT, ValueT]{}
}

func (c *Cache[KeyT, ValueT]) Get(key KeyT) (*ValueT, error) {
	if c.cacheEnabled {
		keyStr, err := c.transformKey(key)
		if err != nil {
			return nil, err
		}
		data, err := c.cache.Get(keyStr)
		if err != nil {
			return nil, err
		}
		result, err := c.deserialize(data)
		if err != nil {
			c.logger.Warn(context.Background(), fmt.Sprint("Cannot deserialize data ", err.Error()))
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("cache is disabled")
}

func (c *Cache[KeyT, ValueT]) Set(key KeyT, value ValueT) error {
	if c.cacheEnabled {
		keyStr, err := c.transformKey(key)
		if err != nil {
			return err
		}
		data, err := c.serialize(value)
		if err != nil {
			c.logger.Warn(context.Background(), fmt.Sprint("Cannot serialize data ", err.Error()))
			return err
		}
		err = c.cache.Set(keyStr, data)
		if err != nil {
			c.logger.Warn(context.Background(), fmt.Sprint("Cannot set data to cache ", err.Error()))
			return err
		}
	}
	return nil
}

func (c *Cache[KeyT, ValueT]) Del(key KeyT) error {
	if c.cacheEnabled {
		keyStr, err := c.transformKey(key)
		if err != nil {
			return err
		}
		err = c.cache.Delete(keyStr)
		if err != nil {
			c.logger.Warn(context.Background(), fmt.Sprint("Cannot delete data from cache ", err.Error()))
			return err
		}
	}
	return nil
}

func (c Cache[KeyT, ValueT]) transformKey(key interface{}) (string, error) {
	if val, ok := key.(uuid.UUID); ok {
		return val.String(), nil
	}
	if val, ok := key.(string); ok {
		return val, nil
	}
	if val, ok := key.(fmt.Stringer); ok {
		return val.String(), nil
	}
	if val, ok := key.(int); ok {
		return strconv.Itoa(val), nil
	}

	if hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil); err == nil {
		return strconv.FormatUint(hash, 10), nil
	}
	c.logger.Error(context.Background(), fmt.Sprint("Cannot convert key to string ", key))
	return "", errors.New("key is not of the type string")
}

func (c *Cache[KeyT, ValueT]) serialize(value ValueT) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Cache[KeyT, ValueT]) deserialize(valueBytes []byte) (*ValueT, error) {
	var value ValueT
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}
