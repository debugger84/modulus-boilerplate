package resolver

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/resolver/validator"
	"boilerplate/internal/user/service"
	"context"
)

type MutationResolver struct {
	registration  *service.Registration
	userValidator *validator.UserValidator
}

func NewMutationResolver(registration *service.Registration, userValidator *validator.UserValidator) *MutationResolver {
	return &MutationResolver{registration: registration, userValidator: userValidator}
}

func (r *MutationResolver) Register(ctx context.Context, request model.RegisterRequest) (*model.User, error) {
	err := r.userValidator.RegisterRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	result, errReg := r.registration.Register(
		ctx, service.RegisterRequest{
			Name:  request.Name,
			Email: request.Email,
		},
	)
	if errReg != nil {
		return nil, err
	}
	return &model.User{
		ID:    result.ID.String(),
		Name:  result.Name,
		Email: result.Email,
	}, nil
}
