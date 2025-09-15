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
