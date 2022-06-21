package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"boilerplate/internal/graph/generated"
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/transformer"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, request model.CreatePostRequest) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *postResolver) Author(ctx context.Context, obj *model.Post) (*model.User, error) {
	id, err := uuid.Parse(obj.Author.ID)
	if err != nil {
		return nil, err
	}
	u, err := r.userLoader.Load(id)
	if err != nil {
		return nil, err
	}
	return transformer.TransformUser(u), nil

	//return transformer.TransformUser(
	//	storage.User{
	//		ID:    id,
	//		Name:  "test",
	//		Email: "",
	//	},
	//), nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Posts(ctx context.Context, count int) ([]*model.Post, error) {
	return r.postQuery.Posts(ctx, count)
}

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }
