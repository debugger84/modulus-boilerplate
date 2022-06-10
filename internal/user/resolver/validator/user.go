package validator

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/dto"
	"context"
	graphql "github.com/debugger84/modulus-graphql"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type UserValidator struct {
}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (v *UserValidator) RegisterRequest(ctx context.Context, request model.RegisterRequest) *gqlerror.Error {
	err := validation.ValidateStructWithContext(
		ctx,
		&request,
		validation.Field(
			&request.Name,
			dto.NameRules()...,
		),
		validation.Field(&request.Email, dto.EmailRules()...),
	)

	if err != nil {
		if errorSet, ok := err.(validation.Errors); ok {
			return graphql.FromValidationErr(ctx, validator.AsAppValidationErrors(errorSet))
		}
	}
	return nil
}
