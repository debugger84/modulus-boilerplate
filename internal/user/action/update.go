package action

import (
	"boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/dto"
	"context"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *UpdateRequest) Validate(ctx context.Context) []application.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
		validation.Field(
			&u.Name,
			dto.NameRules()...,
		),
	)

	return validator.AsAppValidationErrors(err)
}

type Update struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
	logger application.Logger
}

func NewUpdateProcessor(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
	logger application.Logger,
) UpdateProcessor {
	return &Update{finder: finder, saver: saver, logger: logger}
}

func (a *Update) Process(ctx context.Context, request *UpdateRequest) application.ActionResponse {
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
	return application.NewSuccessResponse(
		User{
			Id:   request.Id,
			Name: user.Name,
		},
	)
}

func (a *Update) getUser(ctx context.Context, id string) *dto.User {
	query := a.finder.CreateQuery(ctx)
	query.Id(id)
	user := a.finder.OneByQuery(query)

	return user
}
