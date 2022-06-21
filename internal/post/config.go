package post

import (
	"boilerplate/internal/post/action"
	"boilerplate/internal/post/resolver"
	"boilerplate/internal/post/storage"
	transformer2 "boilerplate/internal/post/transformer"
	transformer "boilerplate/internal/post/transformer/rest"
	application "github.com/debugger84/modulus-application"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
)

type ModuleConfig struct {
	container *dig.Container
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return append(
		action.ServiceProviders(),
		[]interface{}{
			action.NewGetPostProcessor,
			action.NewGetPostListProcessor,

			transformer.NewPostTransformer,
			resolver.NewQueryResolver,
			transformer2.NewPostTransformer,

			func(postTransformer *transformer.PostTransformer) action.PostListTransformer {
				return postTransformer
			},

			func(db *pgxpool.Pool) storage.DBTX {
				return db
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},
		}...,
	)
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func (s *ModuleConfig) ModuleRoutes() []application.RouteInfo {
	var genModuleActions *action.ModuleActions
	err := s.container.Invoke(
		func(dep *action.ModuleActions) {
			genModuleActions = dep
		},
	)
	if err != nil {
		panic("cannot instantiate module dependencies" + err.Error())
	}

	return genModuleActions.Routes()
}
