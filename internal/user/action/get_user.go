package action

import (
	actionError "boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
	"net/http"
)

type GetUserRequest struct {
	Id string `json:"id" validate:"required"`
}

type UserResponse struct {
	Id   string
	Name string
}

type GetUserAction struct {
	runner *application.ActionRunner
	finder *dao.UserFinder
}

func NewGetUserAction(runner *application.ActionRunner, finder *dao.UserFinder) *GetUserAction {
	return &GetUserAction{runner: runner, finder: finder}
}

func (a *GetUserAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(w, r, func(ctx context.Context, request any) application.ActionResponse {
		return a.process(ctx, request.(*GetUserRequest))
	}, &GetUserRequest{})
}

func (a *GetUserAction) process(ctx context.Context, request *GetUserRequest) application.ActionResponse {
	user := a.finder.One(ctx, request.Id)
	if user == nil {
		return actionError.UserNotFound(ctx, request.Id)
	}
	var response UserResponse
	response.Id = request.Id
	response.Name = user.Name
	return application.NewSuccessResponse(response)
}
