package user

import (
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/db"
	"boilerplate/internal/user/resolver"
	"boilerplate/internal/user/resolver/validator"
	"boilerplate/internal/user/service"
	application "github.com/debugger84/modulus-application"
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
			action.NewRegisterAction,
			action.NewGetUsersAction,

			action.NewGetUserProcessor,
			action.NewUpdateProcessor,

			db.NewUserFinder,
			db.NewUserSaver,

			NewModuleActions,

			resolver.NewQueryResolver,
			resolver.NewMutationResolver,
			validator.NewUserValidator,

			service.NewRegistration,
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

	var moduleActions *ModuleActions
	err = s.container.Invoke(
		func(dep *ModuleActions) {
			moduleActions = dep
		},
	)
	if err != nil {
		panic("cannot instantiate module dependencies" + err.Error())
	}

	return append(genModuleActions.Routes(), moduleActions.Routes()...)
}
