-- name: GetPaymentChannels :many
SELECT * FROM payment_channels;

-- name: GetPaymentChannelByID :one
SELECT * FROM payment_channels WHERE id = $1;

-- name: GetPaymentChannelByCode :one
SELECT * FROM payment_channels WHERE code = $1;

-- name: CreatePaymentChannel :one
INSERT INTO payment_channels (
  code,
  name,
  gateway_fee_fixed,
  gateway_fee_percentage_bps,
  platform_fee_fixed,
  platform_fee_percentage_bps
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;