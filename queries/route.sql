-- name: CreateRoute :one
INSERT INTO
    routes (
        number,
        name,
        short_name,
        description,
        "type",
        color,
        agency_id,
        active,
        attributes
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateRoute :one
UPDATE
    routes
SET
    number = coalesce(sqlc.narg('number'), number),
    name = coalesce(sqlc.narg('name'), name),
    short_name = coalesce(sqlc.narg('short_name'), short_name),
    description = coalesce(sqlc.narg('description'), description),
    "type" = coalesce(sqlc.narg('type'), "type"),
    color = coalesce(sqlc.narg('color'), color),
    agency_id = coalesce(sqlc.narg('agency_id'), agency_id),
    active = coalesce(sqlc.narg('active'), active),
    attributes = coalesce(sqlc.narg('attributes'), attributes)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetRoute :one
SELECT
    *
FROM
    routes
WHERE
    id = coalesce(sqlc.narg('id'), id)
LIMIT
    1;

-- name: ListRoute :many
SELECT
    *
FROM
    routes
ORDER BY
    id;
