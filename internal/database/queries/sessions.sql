-- name: GetActiveSessionByTokenID :one
SELECT * FROM sessions WHERE refresh_token_id = $1 AND revoked = FALSE;

-- name: RevokeAllSessionsByUserID :exec
UPDATE sessions SET revoked = TRUE WHERE user_id = $1;

-- name: RevokeSessionByID :exec
UPDATE sessions SET revoked = TRUE WHERE id = $1 AND user_id = $2;

-- name: UpsertSession :one
INSERT INTO sessions (
  id,
  user_id,
  device_id,
  refresh_token_id,
  expires_at,
  meta
) VALUES (
  $1, $2, $3, $4, $5, $6
) 
ON CONFLICT (user_id, device_id)
DO UPDATE SET
  refresh_token_id = EXCLUDED.refresh_token_id,
  expires_at = EXCLUDED.expires_at,
  meta = EXCLUDED.meta,
  last_used_at = now(),
  revoked = FALSE
RETURNING *;