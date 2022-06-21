package transformer

import (
	"boilerplate/internal/post/action"
	"boilerplate/internal/post/storage"
	"boilerplate/internal/user/storage/loader"
	transformer "boilerplate/internal/user/transformer/rest"
	"sync"
)

type PostTransformer struct {
	loader *loader.UserLoader
}

func NewPostTransformer(loader *loader.UserLoader) *PostTransformer {
	return &PostTransformer{loader: loader}
}

func (t *PostTransformer) Transform(post storage.Post) *action.Post {
	author, _ := t.loader.Load(post.AuthorID)
	return &action.Post{
		Author: transformer.TransformUser(author),
		Id:     post.ID.String(),
		Title:  post.Title,
	}
}

func (t *PostTransformer) TransformList(posts []storage.Post) action.PostList {
	if len(posts) == 0 {
		return nil
	}
	wg := sync.WaitGroup{}
	wg.Add(len(posts))
	result := make(action.PostList, len(posts))
	for i, post := range posts {
		go func(i int, post storage.Post) {
			result[i] = *t.Transform(post)
			wg.Done()
		}(i, post)
	}
	wg.Wait()
	return result
}
