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
