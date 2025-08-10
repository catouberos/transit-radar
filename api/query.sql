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
        ebms_id,
        operation_time,
        organization,
        ticketing,
        route_type
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (ebms_id) DO
UPDATE
SET
    number = EXCLUDED.number,
    name = EXCLUDED.name,
    operation_time = EXCLUDED.operation_time,
    organization = EXCLUDED.organization,
    ticketing = EXCLUDED.ticketing,
    route_type = EXCLUDED.route_type RETURNING *;

-- name: GetRoute :one
SELECT
    *
FROM
    routes
WHERE
    id = $1
LIMIT
    1;

-- name: ListRoute :many
SELECT
    *
FROM
    routes
ORDER BY
    id;

-- name: CreateStopType :one
INSERT INTO
    stop_types (name)
VALUES
    ($1) RETURNING *;

-- name: GetStopTypeByName :one
SELECT
    *
FROM
    stop_types
WHERE
    name = $1
LIMIT
    1;

-- name: CreateOrUpdateStop :exec
INSERT INTO
    stops(
        code,
        name,
        type_id,
        ebms_id,
        active,
        latitude,
        longitude
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(ebms_id) DO
UPDATE
SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    type_id = EXCLUDED.type_id,
    active = EXCLUDED.active,
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude;

-- name: GetStopByEbmsID :one
SELECT
    *
FROM
    stops
WHERE
    ebms_id = $1
LIMIT
    1;

-- name: ListStop :many
SELECT
    *
FROM
    stops
ORDER BY
    id;

-- name: GetVariantStopByStopID :many
SELECT
    *
FROM
    variants_stops
WHERE
    stop_id = $1;

-- name: GetVariantStopByVariantID :many
SELECT
    *
FROM
    variants_stops
WHERE
    variant_id = $1;

-- name: CreateVariantStop :one
INSERT INTO
    variants_stops(variant_id, stop_id, order_score)
VALUES
    ($1, $2, $3) RETURNING *;

-- name: ListVariant :many
SELECT
    *
FROM
    variants
ORDER BY
    id;

-- name: GetVariant :one
SELECT
    *
FROM
    variants
WHERE
    id = $1
LIMIT
    1;

-- name: GetVariantByRouteEbmsID :one
SELECT
    *
FROM
    public.variants
    LEFT OUTER JOIN routes ON routes.id = variants.route_id
WHERE
    variants.ebms_id = $1
    AND routes.ebms_id = $2
LIMIT
    1;

-- name: ListVariantByRouteId :many
SELECT
    *
FROM
    variants
WHERE
    route_id = $1;

-- name: CreateOrUpdateVariant :one
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
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (is_outbound, route_id) DO
UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    short_name = EXCLUDED.short_name,
    distance = EXCLUDED.distance,
    duration = EXCLUDED.duration,
    start_stop_name = EXCLUDED.start_stop_name,
    end_stop_name = EXCLUDED.end_stop_name RETURNING *;

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
