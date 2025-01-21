-- name: CreatePost :one
INSERT INTO posts(id, feed_id, title, url, description, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForFeedOfUser :many
SELECT posts.id,
       posts.title,
       posts.url,
       posts.description,
       posts.published_at,
       feeds.id as feedId,
       feeds.name as feedName
FROM public.posts
INNER JOIN feeds ON feeds.id = posts.feed_id
WHERE (select feed_follows.feed_id from feed_follows where feed_follows.feed_id = posts.feed_id) = $1
AND (select feed_follows.user_id from feed_follows where feed_follows.user_id = $2) = $2;