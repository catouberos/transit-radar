-- name: CreateVehicle :one
INSERT INTO
    vehicles(license_plate, "type")
VALUES
    ($1, $2) RETURNING *;

-- name: UpdateVehicle :one
UPDATE
    vehicles
SET
    license_plate = coalesce(sqlc.narg('license_plate'), license_plate),
    "type" = coalesce(sqlc.narg('type'), "type")
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetVehicle :one
SELECT
    *
FROM
    vehicles
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND license_plate = coalesce(sqlc.narg('license_plate'), license_plate)
LIMIT
    1;

-- name: ListVehicle :many
SELECT
    *
FROM
    vehicles
ORDER BY
    id;
