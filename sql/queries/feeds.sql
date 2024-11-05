-- name: Feeds :many
SELECT f.url, f.name,u.name, f.last_fetched_at
FROM feeds f
INNER JOIN users u
ON u.id = f.user_id;

-- name: GetFeed :one
SELECT feeds.id, feeds.name 
FROM feeds
WHERE url = $1 LIMIT 1;

-- name: GetNumberOfFeeds :one
SELECT COUNT(distinct f.id) as N_feeds 
FROM feeds f ;
