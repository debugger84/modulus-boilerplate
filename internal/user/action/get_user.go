package action

import (
	actionError "boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/dto"
	"context"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//type GetUserRequest struct {
//	Id string `json:"id" validate:"required"`
//}

func (u *GetUserRequest) Validate(ctx context.Context) []application.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
	)

	return validator.AsAppValidationErrors(err)
}

//type UserResponse struct {
//	Id   string
//	Name string
//}

//type GetUserAction struct {
//	runner    *application.ActionRunner
//	processor IGetUserProcessor
//}
//type IGetUserProcessor interface {
//	Process(ctx context.Context, request *GetUserRequest) application.ActionResponse
//}

type GetUser struct {
	finder *dao.UserFinder
}

func NewGetUserProcessor(finder *dao.UserFinder) GetUserProcessor {
	return &GetUser{finder: finder}
}

//
//func NewGetUserAction(runner *application.ActionRunner, processor IGetUserProcessor) *GetUserAction {
//	return &GetUserAction{runner: runner, processor: processor}
//}
//
//func (a *GetUserAction) Handle(w http.ResponseWriter, r *http.Request) {
//	a.runner.Run(
//		w, r, func(ctx context.Context, request any) application.ActionResponse {
//			return a.processor.Process(ctx, request.(*GetUserRequest))
//		}, &GetUserRequest{},
//	)
//}

func (a *GetUser) Process(ctx context.Context, request *GetUserRequest) application.ActionResponse {
	user := a.finder.One(ctx, request.Id)
	if user == nil {
		return actionError.UserNotFound(ctx, request.Id)
	}
	var response User
	response.Id = request.Id
	response.Name = user.Name
	return application.NewSuccessResponse(response)
}
