package resolver

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/dto"
	"boilerplate/internal/user/service"
	"context"
)

type MutationResolver struct {
	registration *service.Registration
}

func NewMutationResolver(registration *service.Registration) *MutationResolver {
	return &MutationResolver{registration: registration}
}

func (r *MutationResolver) Register(ctx context.Context, email string, name string) (*model.User, error) {
	user := dto.User{
		Name:  name,
		Email: email,
	}
	result, err := r.registration.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    result.Id,
		Name:  name,
		Email: email,
	}, nil
}
