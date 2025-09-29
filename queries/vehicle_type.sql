-- name: CreateVehicleType :one
INSERT INTO
    vehicle_types(name)
VALUES
    ($1) RETURNING *;

-- name: UpdateVehicleType :one
UPDATE
    vehicle_types
SET
    name = coalesce(sqlc.narg('name'), name)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetVehicleType :one
SELECT
    *
FROM
    vehicle_types
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND name = coalesce(sqlc.narg('name'), name)
LIMIT
    1;

-- name: ListVehicleType :many
SELECT
    *
FROM
    vehicle_types
ORDER BY
    id;
