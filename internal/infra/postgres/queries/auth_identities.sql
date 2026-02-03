-- name: CreateAuthIdentityLocal :exec
INSERT INTO auth_identities (
  id,
  user_id,
  provider,
  provider_id,
  password_hash
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: GetAuthIdentityByEmail :one
SELECT * FROM auth_identities WHERE provider = 'local' AND provider_id = $1;