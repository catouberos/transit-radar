-- name: CreateGeolocation :one
INSERT INTO
    geolocations (
        degree,
        location,
        speed,
        vehicle_id,
        variant_id,
        timestamp
    )
VALUES
    (
        @degree,
        @location :: EWKB,
        @speed,
        @vehicle_id,
        @variant_id,
        @timestamp
    ) RETURNING *;

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
