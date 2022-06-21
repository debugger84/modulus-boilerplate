package action

import (
	"boilerplate/internal/post/storage"
	"context"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *GetPostListRequest) Validate(ctx context.Context) []application.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Count,
			validation.Required.Error("Count is required."),
			validation.Min(1).Error("Count should be positive."),
			validation.Max(100).Error("Count should be less than 100."),
		),
	)

	return validator.AsAppValidationErrors(err)
}

type GetPostList struct {
	finder *storage.Queries
	trans  PostListTransformer
}

type PostListTransformer interface {
	TransformList(posts []storage.Post) PostList
}

func NewGetPostListProcessor(
	finder *storage.Queries,
	trans PostListTransformer,
) GetPostListProcessor {
	return &GetPostList{finder: finder, trans: trans}
}

func (a *GetPostList) Process(ctx context.Context, request *GetPostListRequest) application.ActionResponse {
	posts, err := a.finder.GetPostsFirstPage(ctx, int32(request.Count))
	if err != nil {
		return application.NewServerErrorResponse(ctx, DbError, err)
	}
	return application.NewSuccessResponse(a.trans.TransformList(posts))
}
