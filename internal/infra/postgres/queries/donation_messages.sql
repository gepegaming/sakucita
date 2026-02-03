-- name: CreateDonationMessage :one
INSERT INTO donation_messages (
  id,
  payee_user_id,
  payer_user_id,
  payer_name,
  email,
  message,
  media_type,
  media_url,
  media_start_seconds,
  charged_seconds,
  price_per_second,
  played_at,
  amount,
  currency,
  meta
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;