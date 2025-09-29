-- name: CreateStop :one
INSERT INTO
    stops (
        code,
        name,
        type_id,
        ebms_id,
        active,
        latitude,
        longitude
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateStop :one
UPDATE
    stops
SET
    code = coalesce(sqlc.narg('code'), code),
    name = coalesce(sqlc.narg('name'), name),
    type_id = coalesce(sqlc.narg('type_id'), type_id),
    ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id),
    active = coalesce(sqlc.narg('active'), active),
    latitude = coalesce(sqlc.narg('latitude'), latitude),
    longitude = coalesce(sqlc.narg('longitude'), longitude)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: GetStop :one
SELECT
    *
FROM
    stops
WHERE
    id = coalesce(sqlc.narg('id'), id)
    AND ebms_id = coalesce(sqlc.narg('ebms_id'), ebms_id)
LIMIT
    1;

-- name: ListStop :many
SELECT
    *
FROM
    stops
ORDER BY
    id;
