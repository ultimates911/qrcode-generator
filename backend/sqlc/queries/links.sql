-- name: CreateLink :one
INSERT INTO links (
  original_url,
  hash,
  user_id
) VALUES (
  $1, $2, $3
)
RETURNING id, original_url, hash, created_at, updated_at, user_id;

-- name: GetLinkByHash :one
SELECT id, original_url, hash, created_at, updated_at, user_id FROM links
WHERE hash = $1 LIMIT 1;

-- name: GetLinksByUserID :many
SELECT id, original_url FROM links
WHERE user_id = $1;

-- name: CreateQRCode :one
INSERT INTO qr_codes (
  link_id,
  color,
  background,
  smoothing
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetLinkAndQRCodeByID :one
SELECT
    l.id,
    l.original_url,
    l.hash,
    l.created_at,
    l.updated_at,
    qc.color,
    qc.background,
    qc.smoothing
FROM
    links l
JOIN
    qr_codes qc ON l.id = qc.link_id
WHERE
    l.id = $1 AND l.user_id = $2;