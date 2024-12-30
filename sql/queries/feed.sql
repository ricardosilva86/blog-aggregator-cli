-- name: CreateFeed :one
insert into feeds (id, name, url, created_at, updated_at, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: ListFeeds :many
select * from feeds
join users
on users.id = feeds.user_id;