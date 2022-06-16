package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"boilerplate/internal/graph/model"
	"context"
)

func (r *mutationResolver) Register(ctx context.Context, request model.RegisterRequest) (*model.User, error) {
	return r.userMutation.Register(ctx, request)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.userQuery.User(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context, first int, after *string) (*model.UserList, error) {
	return r.userQuery.Users(ctx, first, after)
}
