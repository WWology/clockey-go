-- name: UpdateScoreboardForGame :exec
INSERT INTO
    public.scoreboards (member, score, game)
VALUES
    ($1, 1, $2) ON CONFLICT ON CONSTRAINT scoreboards_pkey DO
UPDATE
SET
    score = scoreboards.score + 1;

-- name: ShowScoreboardForGame :many
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            score DESC
    ) AS position,
    member,
    score
FROM
    public.scoreboards
WHERE
    game = $1;

-- name: GetMemberScoreForGame :one
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            score DESC
    ) position,
    member,
    score
FROM
    public.scoreboards
WHERE
    game = $1
    AND member = $2;

-- name: GetWinnerForGame :many
SELECT
    *
FROM (
    SELECT
        DENSE_RANK() OVER (
            ORDER BY
                score DESC
        ) position,
        member,
        score
    FROM
        public.scoreboards
    WHERE
        game = $1
)
WHERE
    position = 1;

-- name: ShowGlobalScoreboard :many
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            sum(score) DESC
    ) AS position,
    member,
    sum(score) AS score
FROM
    public.scoreboards
GROUP BY
    member;

-- name: GetMemberGlobalScore :one
SELECT
    DENSE_RANK() OVER (
        ORDER BY
            sum(score) DESC
    ) AS position,
    member,
    sum(score) AS score
FROM
    public.scoreboards
GROUP BY
    member
HAVING
    member = $1;

-- name: GetGlobalWinner :many
WITH
    GlobalRankedLeaderboard AS (
        SELECT
            DENSE_RANK() OVER (
                ORDER BY
                    sum(score) DESC
            ) AS position,
            member,
            sum(score) AS score
        FROM
            public.scoreboards
        GROUP BY
            member
    )
SELECT
    *
FROM
    GlobalRankedLeaderboard
WHERE
    position = 1;

-- name: ClearScoreboard :exec
TRUNCATE TABLE public.scoreboards;
