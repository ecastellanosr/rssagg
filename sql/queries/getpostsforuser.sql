-- name: GetPostsFromUser :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff
ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.updated_at ASC
LIMIT $2;

-- name: GetPostsFromUser1Feed :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff
ON ff.feed_id = p.feed_id
INNER JOIN feeds f
ON f.id = p.feed_id
WHERE ff.user_id = $1 AND f.url = ANY($2::text[])
ORDER BY p.updated_at ASC
LIMIT $3;