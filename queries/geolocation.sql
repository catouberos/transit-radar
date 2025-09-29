-- name: CreateGeolocation :one
INSERT INTO
    geolocations (
        degree,
        latitude,
        longitude,
        speed,
        vehicle_id,
        variant_id,
        "timestamp"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetGeolocation :one
SELECT
    *
FROM
    geolocations
WHERE
    vehicle_id = coalesce(sqlc.narg('vehicle_id'), vehicle_id)
ORDER BY
    "timestamp" DESC
LIMIT
    1;

-- name: ListGeolocation :many
SELECT
    *
FROM
    geolocations
WHERE
    vehicle_id = coalesce(sqlc.narg('vehicle_id'), vehicle_id)
ORDER BY
    "timestamp" DESC
LIMIT
    sqlc.arg('limit');
