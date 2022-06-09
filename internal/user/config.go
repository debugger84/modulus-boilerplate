package user

import (
	action2 "boilerplate/internal/user/action"
	"boilerplate/internal/user/dao"
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
	return []interface{}{
		action2.NewRegisterAction,
		action2.NewGetUserAction,
		action2.NewGetUsersAction,
		action2.NewUpdateAction,

		dao.NewUserFinder,
		dao.NewUserSaver,

		NewModuleActions,

		resolver.NewQueryResolver,
		resolver.NewMutationResolver,
		validator.NewUserValidator,

		service.NewRegistration,
	}
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func (s *ModuleConfig) ModuleRoutes() []application.RouteInfo {
	var moduleActions *ModuleActions
	err := s.container.Invoke(
		func(dep *ModuleActions) {
			moduleActions = dep
		},
	)
	if err != nil {
		panic("cannot instantiate module dependencies" + err.Error())
	}
	return moduleActions.Routes()
}
