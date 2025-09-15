-- name: CreateRoute :one
INSERT INTO
    routes (
        number,
        name,
        ebms_id,
        operation_time,
        organization,
        ticketing,
        route_type
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateRoute :one
UPDATE
    routes
SET
    number = coalesce(sqlc.narg('number'), number),
    name = coalesce(sqlc.narg('name'), name),
    ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id),
    operation_time = coalesce(sqlc.narg('operation_time'), operation_time),
    organization = coalesce(sqlc.narg('organization'), organization),
    ticketing = coalesce(sqlc.narg('ticketing'), ticketing),
    route_type = coalesce(sqlc.narg('route_type'), route_type)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetRoute :one
SELECT
    *
FROM
    routes
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id)
LIMIT
    1;

-- name: ListRoute :many
SELECT
    *
FROM
    routes
ORDER BY
    id;
