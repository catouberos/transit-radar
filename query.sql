-- name: GetLatestByRoute :one
SELECT
    *
FROM
    geolocations
WHERE
    route_id = $1
LIMIT
    1;

-- name: CreateGeolocation :one
INSERT INTO
    geolocations (
        degree,
        latitude,
        longitude,
        speed,
        vehicle_id,
        route_id,
        "timestamp"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;
