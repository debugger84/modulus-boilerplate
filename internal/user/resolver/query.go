package resolver

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/db"
	"context"
	"errors"
)

type QueryResolver struct {
	finder *db.UserFinder
}

func NewQueryResolver(finder *db.UserFinder) *QueryResolver {
	return &QueryResolver{finder: finder}
}

func (r *QueryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user := r.finder.One(ctx, id)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &model.User{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
