// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreatePostRequest struct {
	Title string `json:"title"`
}

type Post struct {
	ID           string  `json:"id"`
	Author       *User   `json:"author"`
	Title        string  `json:"title"`
	Content      *string `json:"content"`
	PreviewImage *string `json:"previewImage"`
}

type PostEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type PostList struct {
	Edges       []*PostEdge `json:"edges"`
	EndCursor   string      `json:"endCursor"`
	HasNextPage bool        `json:"hasNextPage"`
}

type RegisterRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserEdge struct {
	Cursor string `json:"cursor"`
	Node   *User  `json:"node"`
}

type UserList struct {
	Edges       []*UserEdge `json:"edges"`
	EndCursor   string      `json:"endCursor"`
	HasNextPage bool        `json:"hasNextPage"`
}
