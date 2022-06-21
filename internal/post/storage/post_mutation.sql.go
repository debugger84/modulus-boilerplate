// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: post_mutation.sql

package storage

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into "post"."post" (id, author_id, title, previewimage, content, created_at)
values ($1, $2, $3, null, null, now()) RETURNING id, author_id, title, previewimage, content, created_at
`

type CreateUserParams struct {
	ID       uuid.UUID `db:"id" json:"id"`
	AuthorID uuid.UUID `db:"author_id" json:"authorID"`
	Title    string    `db:"title" json:"title"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Post, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.AuthorID, arg.Title)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Previewimage,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}