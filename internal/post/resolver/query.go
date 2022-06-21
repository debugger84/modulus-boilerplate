package resolver

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/post/storage"
	"boilerplate/internal/post/transformer"
	"context"
	"errors"
)

type QueryResolver struct {
	finder *storage.Queries
	trans  *transformer.PostTransformer
}

func NewQueryResolver(finder *storage.Queries, trans *transformer.PostTransformer) *QueryResolver {
	return &QueryResolver{finder: finder, trans: trans}
}

func (r *QueryResolver) Posts(ctx context.Context, count int) ([]*model.Post, error) {
	posts, err := r.finder.GetPostsFirstPage(ctx, int32(count))
	if err != nil {
		return nil, errors.New("db error")
	}

	return r.trans.TransformList(posts), nil
}
