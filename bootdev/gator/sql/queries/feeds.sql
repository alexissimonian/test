-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds
WHERE url = $1;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name AS username
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS(
INSERT INTO feed_follows(
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
)
SELECT inserted_feed_follow.*, users.name AS username, feeds.name AS feedname
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT f.id, f.created_at, f.updated_at, f.user_id, f.feed_id, users.name AS username, feeds.name AS feedname
FROM feed_follows f
INNER JOIN users ON users.id = f.user_id
INNER JOIN feeds ON feeds.id = f.feed_id
WHERE f.user_id = $1;

-- name: RemoveFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
