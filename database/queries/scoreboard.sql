-- name: UpdateScoreboardForGame :execrows
INSERT INTO
    scoreboards (member, score, game)
VALUES
    ($1, $2, $3) ON CONFLICT (member) DO
UPDATE
SET
    score = score + 1;

-- name: ShowScoreboardForGame :many
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            score DESC
    ) position,
    member,
    score
FROM
    scoreboards
WHERE
    game = $1;

-- name: ClearScoreboardForGame :execrows
DELETE FROM
    scoreboards
WHERE
    game = $1;

-- name: GetWinnerForGame :many
SELECT
    position, member, score
FROM (
    SELECT
        DENSE_RANK() OVER (
            ORDER BY
                score DESC
        ) position,
        member,
        score
    FROM
        scoreboards
    WHERE
        game = $1
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
        game = $1
)
WHERE
    member = $2;

-- name: ShowGlobalScoreboard :many
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            score DESC
    ) position,
    member,
    sum(score) score
FROM
    scoreboards
GROUP BY
    member;

-- name: GetMemberGlobalScore :one
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            sum(score) DESC
    ) position,
    sum(score) score
FROM
    scoreboards
WHERE
    member = $1;
