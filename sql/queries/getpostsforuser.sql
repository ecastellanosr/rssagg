-- name: GetPostsFromUser :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff
ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.updated_at ASC
LIMIT $2;