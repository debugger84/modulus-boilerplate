package pgx

import (
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
)

type ModuleConfig struct {
	container   *dig.Container
	PgDsn       string
	SlowQueryMs int
}

func (s *ModuleConfig) InitConfig(config application.Config) error {
	if s.PgDsn == "" {
		s.PgDsn = config.GetEnv("PG_DSN")
	}
	if s.SlowQueryMs == 0 {
		s.SlowQueryMs = config.GetEnvAsInt("DB_SLOW_QUERY_LOGGING_LIMIT")
	}

	return nil
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{
		NewPgxPool,
		NewPgxLogger,
		func() *ModuleConfig {
			return s
		},
	}
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func NewPgxPool(cfg *ModuleConfig, logger *PgxLogger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(cfg.PgDsn)
	if err != nil {
		panic("cannot parse pg dsn" + err.Error())
	}
	config.ConnConfig.Logger = logger
	//config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	conn.ConnInfo().RegisterDataType(
	//		pgtype.DataType{
	//			Value: &pgtypeuuid.UUID{},
	//			Name:  "uuid",
	//			OID:   pgtype.UUIDOID,
	//		},
	//	)
	//	return nil
	//}
	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic("cannot establish connection" + err.Error())
	}

	return dbPool
}
