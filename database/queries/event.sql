-- name: CreateEvent :exec
INSERT INTO events (name, time, type, gardener, hours) VALUES (?, ?, ?, ?, ?);

-- name: DeleteEvent :exec
DELETE FROM events
WHERE name = ? AND time = ?;

-- name: GetEventsForGardener :many
SELECT
    *
FROM
    events
WHERE time BETWEEN @start AND @end
AND gardener = ?;

-- name: GetEventsForGame :many
SELECT
    *
FROM
    events
WHERE time BETWEEN @start AND @end
AND type = ?;
