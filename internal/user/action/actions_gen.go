// Package action provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/debugger84/oapi-codegen version v1.11.5 DO NOT EDIT.
package action

import (
	"context"
	"net/http"

	application "github.com/debugger84/modulus-application"
)

// User defines model for User.
type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// UserId defines model for userId.
type UserId = string

// UpdateUserBody defines model for UpdateUserBody.
type UpdateUserBody struct {
	Name string `json:"name"`
}

// UpdateJSONRequestBody defines body for Update for application/json ContentType.
type UpdateJSONRequestBody UpdateUserBody

// ------------- Code generation for the "GetUser" -------------
type GetUserRequest struct {
	// Path parameter "id"
	Id UserId `json:"id"`
}
type GetUserAction struct {
	runner    *application.ActionRunner
	processor GetUserProcessor
}
type GetUserProcessor interface {
	Process(ctx context.Context, request *GetUserRequest) application.ActionResponse
}

func NewGetUserAction(runner *application.ActionRunner, processor GetUserProcessor) *GetUserAction {
	return &GetUserAction{runner: runner, processor: processor}
}

func (a *GetUserAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(
		w, r, func(ctx context.Context, request any) application.ActionResponse {
			return a.processor.Process(ctx, request.(*GetUserRequest))
		}, &GetUserRequest{},
	)
}

// ------------- Code generation for the "Update" -------------
type UpdateRequest struct {
	// Path parameter "id"
	Id string `json:"id"`
	UpdateUserBody
}
type UpdateAction struct {
	runner    *application.ActionRunner
	processor UpdateProcessor
}
type UpdateProcessor interface {
	Process(ctx context.Context, request *UpdateRequest) application.ActionResponse
}

func NewUpdateAction(runner *application.ActionRunner, processor UpdateProcessor) *UpdateAction {
	return &UpdateAction{runner: runner, processor: processor}
}

func (a *UpdateAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(
		w, r, func(ctx context.Context, request any) application.ActionResponse {
			return a.processor.Process(ctx, request.(*UpdateRequest))
		}, &UpdateRequest{},
	)
}

//------ services initialization
type ModuleActions struct {
	routes *application.Routes
}

func NewModuleActions(
	getUserAction *GetUserAction,
	updateAction *UpdateAction,
) *ModuleActions {
	routes := application.NewRoutes()

	routes.Get(
		"/users/:id",
		getUserAction.Handle,
	)

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

func ServiceProviders() []interface{} {
	return []interface{}{
		NewModuleActions,

		NewGetUserAction,
		NewUpdateAction,
	}
}