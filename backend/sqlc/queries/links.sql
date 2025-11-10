-- name: CreateLink :one
INSERT INTO links (
  original_url,
  hash,
  user_id,
  name
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, original_url, hash, created_at, updated_at, user_id, name;

-- name: GetLinkByHash :one
SELECT id, original_url, hash, created_at, updated_at, user_id, name FROM links
WHERE hash = $1 LIMIT 1;

-- name: GetLinksByUserID :many
SELECT id, original_url, name FROM links
WHERE user_id = $1;

-- name: SearchLinksByName :many
SELECT id, original_url, name FROM links
WHERE user_id = $1 AND name ILIKE '%' || $2 || '%';

-- name: GetLinksSummaryByUser :many
SELECT 
  l.id,
  l.original_url,
  l.name,
  l.created_at,
  COALESCE(COUNT(t.id), 0) AS transitions_count
FROM links l
LEFT JOIN transitions t ON t.link_id = l.id
WHERE l.user_id = $1
GROUP BY l.id
ORDER BY l.created_at DESC;

-- name: SearchLinksSummaryByName :many
SELECT 
  l.id,
  l.original_url,
  l.name,
  l.created_at,
  COALESCE(COUNT(t.id), 0) AS transitions_count
FROM links l
LEFT JOIN transitions t ON t.link_id = l.id
WHERE l.user_id = $1 AND l.name ILIKE '%' || $2 || '%'
GROUP BY l.id
ORDER BY l.created_at DESC;

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
    l.name,
    qc.color,
    qc.background,
    qc.smoothing
FROM
    links l
JOIN
    qr_codes qc ON l.id = qc.link_id
WHERE
    l.id = $1 AND l.user_id = $2;

-- name: UpdateLinkURL :execrows
UPDATE links
SET
    original_url = $1,
    updated_at = now()
WHERE
    id = $2 AND user_id = $3;

-- name: UpdateQRCodeParams :exec
UPDATE qr_codes
SET
    color = $1,
    background = $2,
    smoothing = $3
WHERE
    link_id = $4;

-- name: CreateTransition :exec
INSERT INTO transitions (
  link_id,
  country,
  city,
  referer,
  user_agent,
  browser,
  os
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: GetTransitionsByLinkID :many
SELECT
  t.id,
  t.country,
  t.city,
  t.referer,
  t.user_agent,
  t.browser,
  t.os,
  t.created_at
FROM transitions t
JOIN links l ON l.id = t.link_id
WHERE t.link_id = $1 AND l.user_id = $2
ORDER BY t.created_at DESC;

-- name: DeleteTransitionsByLinkID :exec
DELETE FROM transitions WHERE link_id = $1;

-- name: DeleteQRCodeByLinkID :exec
DELETE FROM qr_codes WHERE link_id = $1;

-- name: DeleteLink :execrows
DELETE FROM links WHERE id = $1 AND user_id = $2;