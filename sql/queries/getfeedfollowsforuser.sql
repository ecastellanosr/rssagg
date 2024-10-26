-- name: GetFeedFollowsforUser :many
SELECT ff.*, u.name AS Username, f.name AS feed_name FROM feed_follows ff
INNER JOIN users u
ON u.id = ff.user_id
INNER JOIN feeds f
ON f.id = ff.feed_id
WHERE ff.user_id = $1;

-- name: getfeedfollows :many
SELECT ff.*, u.name AS Username, f.name AS Feed_name FROM feed_follows ff
INNER JOIN users u
ON u.id = ff.user_id
INNER JOIN feeds f
ON f.id = ff.feed_id
WHERE ff.user_id = $1;