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

type GetUser struct {
	finder *dao.UserFinder
}

func NewGetUserProcessor(finder *dao.UserFinder) GetUserProcessor {
	return &GetUser{finder: finder}
}

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
