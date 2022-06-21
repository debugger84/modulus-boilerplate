
-- name: CreateUser :one
insert into "post"."post" (id, author_id, title, previewimage, content, created_at)
values ($1, $2, $3, null, null, now()) RETURNING *;