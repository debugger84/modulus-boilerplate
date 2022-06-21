package transformer

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/post/storage"
)

type PostTransformer struct {
}

func NewPostTransformer() *PostTransformer {
	return &PostTransformer{}
}

func (t *PostTransformer) Transform(post storage.Post) *model.Post {
	return &model.Post{
		ID:           post.ID.String(),
		Author:       &model.User{ID: post.AuthorID.String()},
		Title:        post.Title,
		Content:      nil,
		PreviewImage: &post.Previewimage.String,
	}
}

func (t *PostTransformer) TransformList(posts []storage.Post) []*model.Post {
	result := make([]*model.Post, len(posts))
	for i, post := range posts {
		result[i] = t.Transform(post)
	}
	return result
}
