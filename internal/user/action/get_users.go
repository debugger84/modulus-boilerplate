package action

import (
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
	"net/http"
)

type GetUsersRequest struct {
	Count int `json:"count" validate:"required,gte=0,lte=10"`
}

type UsersResponse struct {
	List []UserResponse `json:"list"`
}

type GetUsersAction struct {
	runner *application.ActionRunner
	finder *dao.UserFinder
}

func NewGetUsersAction(runner *application.ActionRunner, finder *dao.UserFinder) *GetUsersAction {
	return &GetUsersAction{runner: runner, finder: finder}
}

func (a *GetUsersAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(w, r, func(ctx context.Context, request any) application.ActionResponse {
		return a.process(ctx, request.(*GetUsersRequest))
	}, &GetUsersRequest{})
}

func (a *GetUsersAction) process(ctx context.Context, request *GetUsersRequest) application.ActionResponse {
	query := a.finder.CreateQuery(ctx)
	query.NewerFirst()
	users := a.finder.ListByQuery(query, request.Count)
	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{Id: user.Id, Name: user.Name}
	}

	return application.NewSuccessResponse(UsersResponse{
		List: response,
	})
}
