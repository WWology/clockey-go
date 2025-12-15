-- name: CreateEvent :exec
INSERT INTO events (name, time, type, gardener, hours) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteEvent :exec
DELETE FROM events
WHERE name = $1 AND time = $2 AND type = $3 AND hours = $4;

-- name: GetEventsForGardener :many
SELECT
    *
FROM
    events
WHERE time BETWEEN @start_time AND @end_time
AND gardener = $1;

-- name: GetEventsForGame :many
SELECT
    *
FROM
    events
WHERE time BETWEEN @start_time AND @end_time
AND type = $1;
