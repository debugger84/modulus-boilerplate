package errors

import (
	"context"
	"errors"
	"fmt"
	application "github.com/debugger84/modulus-application"
)

const postNotFound application.ErrorIdentifier = "PostNotFound"
const cannotUpdatePost application.ErrorIdentifier = "CannotUpdatePost"

func PostNotFound(ctx context.Context, id string) application.ActionResponse {
	return application.ActionResponse{
		StatusCode: 404,
		Error: &application.ActionError{
			Ctx:              ctx,
			Identifier:       postNotFound,
			Err:              errors.New(fmt.Sprintf("Post with id %s is not found", id)),
			ValidationErrors: nil,
		},
	}
}

func CannotUpdatePost(ctx context.Context, id string) application.ActionResponse {
	return application.ActionResponse{
		StatusCode: 422,
		Error: &application.ActionError{
			Ctx:              ctx,
			Identifier:       cannotUpdatePost,
			Err:              errors.New(fmt.Sprintf("Post with id %s cannot be updated", id)),
			ValidationErrors: nil,
		},
	}
}
