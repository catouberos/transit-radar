-- WARN legacy queries, will be migrated away
--
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
