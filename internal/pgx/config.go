package pgx

import (
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
)

type ModuleConfig struct {
	container *dig.Container
	PgDsn     string
}

func (s *ModuleConfig) InitConfig(config application.Config) error {
	if s.PgDsn == "" {
		s.PgDsn = config.GetEnv("PG_DSN")
	}
	return nil
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{
		NewPgxPool,
		func() *ModuleConfig {
			return s
		},
	}
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func NewPgxPool(cfg *ModuleConfig) *pgxpool.Pool {
	dbPool, err := pgxpool.Connect(context.Background(), cfg.PgDsn)

	if err != nil {
		panic("cannot establish connection" + err.Error())
	}

	return dbPool
}
