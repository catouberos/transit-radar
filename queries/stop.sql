-- name: CreateStop :one
INSERT INTO
    stops (
        parent_id,
        code,
        name,
        "type",
        active,
        location,
        attributes
    )
VALUES
    (
        @parent_id,
        @code,
        @name,
        @type,
        @active,
        sqlc.arg('location') :: EWKB,
        @attributes
    ) RETURNING *;

-- name: UpdateStop :one
UPDATE
    stops
SET
    parent_id = coalesce(sqlc.narg('parent_id'), parent_id),
    code = coalesce(sqlc.narg('code'), code),
    name = coalesce(sqlc.narg('name'), name),
    "type" = coalesce(sqlc.narg('type'), "type"),
    active = coalesce(sqlc.narg('active'), active),
    location = coalesce(sqlc.narg('location'), location),
    attributes = coalesce(sqlc.narg('attributes'), attributes)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetStop :one
SELECT
    *
FROM
    stops
WHERE
    id = coalesce(sqlc.narg('id'), id)
LIMIT
    1;

-- name: ListStop :many
SELECT
    *
FROM
    stops
ORDER BY
    id;
