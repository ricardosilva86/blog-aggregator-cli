-- name: CreateFeed :one
insert into feeds (id, name, url, created_at, updated_at, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: ListFeeds :many
select * from feeds
join users
on users.id = feeds.user_id
where user_id = $1;

-- name: GetFeedByURL :one
select * from feeds
where url = $1;

-- name: MarkFeedFetched :one
update feeds set updated_at = now(), last_fetched_at = now()
where id = $1
returning *;

-- name: GetNextFeedToFetch :one
select * from feeds
order by feeds.last_fetched_at NULLS FIRST
limit 1;