-- name: UpdateScoreboardForGame :execrows
INSERT INTO
    scoreboards (member, score, game)
VALUES
    (?, ?, ?) ON CONFLICT (member) DO
UPDATE
SET
    score = score + 1;

-- name: ShowScoreboardForGame :many
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            score DESC
    ) position,
    id,
    score
FROM
    scoreboards
WHERE
    game = ?;

-- name: ClearScoreboardForGame :execrows
DELETE FROM
    scoreboards
WHERE
    game = ?;

-- name: GetWinnerForGame :many
SELECT
    position, id, score
FROM (
    SELECT
        DENSE_RANK() OVER (
            ORDER BY
                score DESC
        ) position,
        id,
        score
    FROM
        scoreboards
    WHERE
        game = ?
)
WHERE
    position = 1;

-- name: GetMemberScoreForGame :one
SELECT
    position, score
FROM (
    SELECT
        DENSE_RANK() OVER (
            ORDER BY
                score DESC
        ) position,
        score
    FROM
        scoreboards
    WHERE
        game = ?
)
WHERE
    id = ?;
