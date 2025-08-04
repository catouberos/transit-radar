-- name: GetRouteByVariantID :one
SELECT
    *
FROM
    geolocations
WHERE
    variant_id = $1
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
        variant_id,
        "timestamp"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CreateOrUpdateRoute :one
INSERT INTO
    routes (
        number,
        name,
        ebms_id
    )
VALUES
    ($1, $2, $3) ON CONFLICT (ebms_id) DO
UPDATE
SET
    number = EXCLUDED.number,
    name = EXCLUDED.name RETURNING *;

-- name: CreateOrUpdateVariant :one
INSERT INTO
    variants (
        name,
        ebms_id,
        is_outbound,
        route_id
    )
VALUES
    ($1, $2, $3, $4) ON CONFLICT (is_outbound, route_id) DO
UPDATE
SET
    name = EXCLUDED.name RETURNING *;

-- name: GetVariantByRouteIDAndOutbound :one
SELECT
    *
FROM
    variants
WHERE
    route_id = $1
    AND is_outbound = $2
LIMIT
    1;

-- name: GetRouteByEbmsID :one
SELECT
    *
FROM
    routes
WHERE
    ebms_id = $1
LIMIT
    1;

-- name: CreateVehicle :one
INSERT INTO
    vehicles(license_plate)
VALUES
    ($1) RETURNING *;

-- name: GetVehicleByLicensePlate :one
SELECT
    *
FROM
    vehicles
WHERE
    license_plate = $1
LIMIT
    1;
