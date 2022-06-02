package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"boilerplate/internal/graph/model"
	"context"
)

func (r *mutationResolver) Register(ctx context.Context, email string, name string) (*model.User, error) {
	return r.userMutation.Register(ctx, email, name)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.userQuery.User(ctx, id)
}
