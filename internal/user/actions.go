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
	getUserAction *action2.GetUserAction,
	getUsersAction *action2.GetUsersAction,
	updateAction *action2.UpdateAction,
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
	routes.Get("/users/:id", getUserAction.Handle)

	routes.Put(
		"/users/:id",
		updateAction.Handle,
	)

	return &ModuleActions{
		routes: routes,
	}
}

func (a *ModuleActions) Routes() []application.RouteInfo {
	return a.routes.GetRoutesInfo()
}
