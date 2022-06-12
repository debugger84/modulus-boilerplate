package user

import (
	action2 "boilerplate/internal/user/action"
	application "github.com/debugger84/modulus-application"
)

type ModuleActions struct {
	routes *application.Routes
}

func NewModuleActions(
	registerAction *action2.RegisterAction,
	getUsersAction *action2.GetUsersAction,
) *ModuleActions {
	routes := application.NewRoutes()
	routes.Post(
		"/users",
		registerAction.Handle,
	)
	routes.Get(
		"/users",
		getUsersAction.Handle,
	)

	return &ModuleActions{
		routes: routes,
	}
}

func (a *ModuleActions) Routes() []application.RouteInfo {
	return a.routes.GetRoutesInfo()
}
