package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	userResolver "boilerplate/internal/user/resolver"
	application "github.com/debugger84/modulus-application"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger       application.Logger
	userQuery    *userResolver.QueryResolver
	userMutation *userResolver.MutationResolver
}

func NewResolver(
	logger application.Logger,
	userQuery *userResolver.QueryResolver,
	userMutation *userResolver.MutationResolver,
) *Resolver {
	return &Resolver{
		userQuery:    userQuery,
		userMutation: userMutation,
		logger:       logger,
	}
}
