-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (feed_id, user_id)
    VALUES ($1, $2)
    RETURNING *
)

SELECT
    inserted_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follows
JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
JOIN users ON inserted_feed_follows.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*,
    users.name AS user_name, 
    users.id AS user_id, 
    feed_follows.feed_id AS feed_name
FROM feed_follows 
JOIN users ON feed_follows.user_id = users.id
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows 
WHERE user_id = $1 
    AND feed_id = $2;