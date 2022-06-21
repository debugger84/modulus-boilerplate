-- name: GetPost :one
select * from "post"."post" where id = $1 LIMIT 1;

-- name: GetNewerPosts :many
select * from "post"."post" order by "post".created_at DESC LIMIT $1;

-- name: GetPostsFirstPage :many
select * from "post"."post" order by "post".created_at DESC, "post".id DESC LIMIT $1;

-- name: GetPostsAfterCursor :many
select * from "post"."post"
where
    created_at < $1 OR
    (created_at = $2 AND id < $3)
order by "post".created_at DESC, "post".id DESC LIMIT $4;