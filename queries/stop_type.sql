-- name: CreateStopType :one
INSERT INTO
    stop_types (name)
VALUES
    ($1) RETURNING *;

-- name: GetStopTypeByName :one
SELECT
    *
FROM
    stop_types
WHERE
    name = $1
LIMIT
    1;
