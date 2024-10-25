-- name: Feeds :many
SELECT f.url, f.name,u.name 
FROM feeds f
INNER JOIN users u
ON u.id = f.user_id;