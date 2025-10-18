-- name: CreateVariant :one
INSERT INTO
    variants (
        route_id,
        name,
        short_name,
        description,
        distance,
        direction,
        duration,
        attributes
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateVariant :one
UPDATE
    variants
SET
    route_id = coalesce(sqlc.narg('route_id'), route_id),
    name = coalesce(sqlc.narg('name'), name),
    short_name = coalesce(sqlc.narg('short_name'), short_name),
    description = coalesce(sqlc.narg('description'), description),
    distance = coalesce(sqlc.narg('distance'), distance),
    direction = coalesce(sqlc.narg('direction'), direction),
    duration = coalesce(sqlc.narg('duration'), duration),
    attributes = coalesce(sqlc.narg('attributes'), attributes)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetVariant :one
SELECT
    *
FROM
    variants
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND route_id = coalesce(sqlc.narg('route_id'), route_id)
LIMIT
    1;

-- name: ListVariant :many
SELECT
    *
FROM
    variants
WHERE
    route_id = coalesce(sqlc.narg('route_id'), route_id)
ORDER BY
    id;
