package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"boilerplate/internal/graph/generated"
	post "boilerplate/internal/post/resolver"
	user "boilerplate/internal/user/resolver"
	"boilerplate/internal/user/storage/loader"
	application "github.com/debugger84/modulus-application"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger       application.Logger
	userQuery    *user.QueryResolver
	userMutation *user.MutationResolver
	postQuery    *post.QueryResolver
	userLoader   *loader.UserLoader
}

func NewResolver(
	logger application.Logger,
	userQuery *user.QueryResolver,
	userMutation *user.MutationResolver,
	postQuery *post.QueryResolver,
	userLoader *loader.UserLoader,
) *Resolver {
	return &Resolver{
		userQuery:    userQuery,
		userMutation: userMutation,
		logger:       logger,
		postQuery:    postQuery,
		userLoader:   userLoader,
	}
}

func (r Resolver) GetDirectives() generated.DirectiveRoot {
	return generated.DirectiveRoot{}
}
