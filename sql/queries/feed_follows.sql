-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    insert into feed_follows(id, user_id, feed_id, created_at, updated_at)
        values ($1, $2, $3, $4, $5)
    returning *)
SELECT inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;


-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id,
       feeds.id as feedId,
       feeds.name as feedName,
       users.id as userId,
       users.name as userName
FROM public.feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where feed_follows.user_id = $1
and (select id from feeds where url = $2) = feed_id;