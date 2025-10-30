-- name: CreateEvent :exec
INSERT INTO events (name, time, type, gardener, hours) VALUES (?, ?, ?, ?, ?);
