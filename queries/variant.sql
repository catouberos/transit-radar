-- name: CreateVariant :one
INSERT INTO
    variants (
        name,
        ebms_id,
        is_outbound,
        route_id,
        description,
        short_name,
        distance,
        duration,
        start_stop_name,
        end_stop_name
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpdateVariant :one
UPDATE
    variants
SET
    name = coalesce(sqlc.narg('name'), name),
    ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id),
    is_outbound = coalesce(sqlc.narg('is_outbound'), is_outbound),
    route_id = coalesce(sqlc.narg('route_id'), route_id),
    description = coalesce(sqlc.narg('description'), description),
    short_name = coalesce(sqlc.narg('short_name'), short_name),
    distance = coalesce(sqlc.narg('distance'), distance),
    duration = coalesce(sqlc.narg('duration'), duration),
    start_stop_name = coalesce(sqlc.narg('start_stop_name'), start_stop_name),
    end_stop_name = coalesce(sqlc.narg('end_stop_name'), end_stop_name)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetVariant :one
SELECT
    *
FROM
    variants
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id)
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
    AND is_outbound = coalesce(sqlc.narg('is_outbound'), is_outbound)
ORDER BY
    id;
