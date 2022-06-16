// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user_query.sql

package storage

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

const getUser = `-- name: GetUser :one
select id, name, email, registered_at, settings, contacts from "user"."user" where id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.RegisteredAt,
		&i.Settings,
		&i.Contacts,
	)
	return i, err
}

const getUsersAfterCursor = `-- name: GetUsersAfterCursor :many
select id, name, email, registered_at, settings, contacts from "user"."user"
where
    registered_at < $1 OR
    (registered_at = $2 AND id < $3)
order by "user".registered_at DESC, "user".id DESC LIMIT $4
`

type GetUsersAfterCursorParams struct {
	RegisteredAt   time.Time `db:"registered_at" json:"registeredAt"`
	RegisteredAt_2 time.Time `db:"registered_at_2" json:"registeredAt2"`
	ID             uuid.UUID `db:"id" json:"id"`
	Limit          int32     `db:"limit" json:"limit"`
}

func (q *Queries) GetUsersAfterCursor(ctx context.Context, arg GetUsersAfterCursorParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersAfterCursor,
		arg.RegisteredAt,
		arg.RegisteredAt_2,
		arg.ID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.RegisteredAt,
			&i.Settings,
			&i.Contacts,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersFirstPage = `-- name: GetUsersFirstPage :many
select id, name, email, registered_at, settings, contacts from "user"."user" order by "user".registered_at DESC, "user".id DESC LIMIT $1
`

func (q *Queries) GetUsersFirstPage(ctx context.Context, limit int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersFirstPage, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.RegisteredAt,
			&i.Settings,
			&i.Contacts,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
