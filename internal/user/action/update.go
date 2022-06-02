package action

import (
	"boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/dto"
	"context"
	application "github.com/debugger84/modulus-application"
	"net/http"
)

type UpdateRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type UpdateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UpdateAction struct {
	runner *application.ActionRunner
	finder *dao.UserFinder
	saver  *dao.UserSaver
	logger application.Logger
}

func NewUpdateAction(runner *application.ActionRunner, finder *dao.UserFinder, saver *dao.UserSaver, logger application.Logger) *UpdateAction {
	return &UpdateAction{runner: runner, finder: finder, saver: saver, logger: logger}
}

func (a *UpdateAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(w, r, func(ctx context.Context, request any) application.ActionResponse {
		return a.process(ctx, request.(*UpdateRequest))
	}, &UpdateRequest{})
}

func (a *UpdateAction) process(ctx context.Context, request *UpdateRequest) application.ActionResponse {
	user := a.getUser(ctx, request.Id)
	if user == nil {
		return errors.UserNotFound(ctx, request.Id)
	}
	user.Name = request.Name
	err := a.saver.Update(ctx, *user)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return errors.CannotUpdateUser(ctx, request.Id)
	}
	return application.NewSuccessResponse(UpdateResponse{
		Id:   request.Id,
		Name: user.Name,
	})
}

func (a *UpdateAction) getUser(ctx context.Context, id string) *dto.User {
	query := a.finder.CreateQuery(ctx)
	query.Id(id)
	user := a.finder.OneByQuery(query)

	return user
}
