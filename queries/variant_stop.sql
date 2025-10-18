-- name: CreateVariantStop :one
INSERT INTO
    variants_stops(
        variant_id,
        stop_id,
        order_score
    )
VALUES
    ($1, $2, $3) RETURNING *;

-- name: UpdateVariantStop :one
UPDATE
    variants_stops
SET
    stop_id = $3
WHERE
    variant_id = $1
    AND order_score = $2 RETURNING *;

-- name: ListVariantStop :many
SELECT
    *
FROM
    variants_stops
WHERE
    variant_id = coalesce(sqlc.narg('variantID'), variant_id)
    AND stop_id = coalesce(sqlc.narg('stopID'), stop_id)
ORDER BY
    order_score ASC;

-- name: DeleteVariantStop :exec
DELETE FROM
    variants_stops
WHERE
    variant_id = $1
    AND order_score = $2;
