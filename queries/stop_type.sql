-- name: CreateStopType :one
INSERT INTO
    stop_types(name)
VALUES
    ($1) RETURNING *;

-- name: UpdateStopType :one
UPDATE
    stop_types
SET
    name = coalesce(sqlc.narg('name'), name)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetStopType :one
SELECT
    *
FROM
    stop_types
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND name = coalesce(sqlc.narg('name'), name)
LIMIT
    1;

-- name: ListStopType :many
SELECT
    *
FROM
    stop_types
ORDER BY
    id;
