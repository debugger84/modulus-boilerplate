package action

import (
	actionError "boilerplate/internal/post/action/errors"
	"boilerplate/internal/post/storage"
	"boilerplate/internal/user/dto"
	userStorage "boilerplate/internal/user/storage"
	"boilerplate/internal/user/storage/loader"
	transformer "boilerplate/internal/user/transformer/rest"
	"context"
	"errors"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/vikstrous/dataloadgen"
)

const DbError application.ErrorIdentifier = "DbError"

func (u *GetPostRequest) Validate(ctx context.Context) []application.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
	)

	return validator.AsAppValidationErrors(err)
}

type GetPost struct {
	finder  *storage.Queries
	loader  *loader.UserLoader
	loader2 *dataloadgen.Loader[uuid.UUID, userStorage.User]
}

func NewGetPostProcessor(
	finder *storage.Queries,
	loader *loader.UserLoader,
	loader2 *dataloadgen.Loader[uuid.UUID, userStorage.User],
) GetPostProcessor {
	return &GetPost{finder: finder, loader: loader, loader2: loader2}
}

func (a *GetPost) Process(ctx context.Context, request *GetPostRequest) application.ActionResponse {
	id, _ := uuid.Parse(request.Id)
	post, err := a.finder.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return actionError.PostNotFound(ctx, request.Id)
		} else {
			return application.NewServerErrorResponse(ctx, DbError, err)
		}
	}
	var response Post
	response.Id = post.ID.String()
	response.Title = post.Title
	author, _ := a.loader.Load(post.AuthorID)
	response.Author = transformer.TransformUser(author)

	return application.NewSuccessResponse(response)
}
