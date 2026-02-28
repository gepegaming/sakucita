-- name: CreateTransaction :one
INSERT INTO transactions (
  id,
  donation_message_id,
  payment_channel_id,
  payer_user_id,
  payee_user_id,
  gross_paid_amount,
  gateway_fee_fixed,
  gateway_fee_percentage_bps,
  gateway_fee_amount,
  platform_fee_fixed,
  platform_fee_percentage_bps,
  platform_fee_amount,
  fee_fixed,
  fee_percentage_bps,
  fee_amount,
  net_amount,
  currency,
  status,
  external_reference
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetByExternalReference :one
SELECT * FROM transactions WHERE external_reference = $1 LIMIT 1;

-- name: MarkTransactionAs :exec
UPDATE transactions
SET status = $2
WHERE id = $1;


-- name: UpdateTransactionExternalReferenceAndStatus :exec
UPDATE transactions
SET external_reference = $2, status = $3
WHERE id = $1;

-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET status = $2
WHERE id = $1;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1 LIMIT 1;

-- name: GetTransactionByDonationMessageID :one
SELECT * FROM transactions WHERE donation_message_id = $1 LIMIT 1;

-- name: GetTransactionByExternalReference :one
SELECT * FROM transactions WHERE external_reference = $1 LIMIT 1;