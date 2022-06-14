package db

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/gofrs/uuid"
)

type Settings struct {
	Incognito bool `json:"incognito"`
}

func (j Settings) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Settings) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

func (f *UserFinder) One(ctx context.Context, id string) *User {
	uid, err := uuid.FromString(id)
	if err != nil {
		return nil
	}
	query := f.CreateQuery(ctx)
	query.ID(uid)

	return f.OneByQuery(query)
}

func (p *UserQuery) NewerFirst() {
	p.db = p.db.Order(UserTable + ".registered_at DESC")
}
