-- name: CreateDonationMessage :one
INSERT INTO donation_messages (
  id,
  payee_user_id,
  payer_user_id,
  payer_name,
  message,
  email,
  media_type,
  tts_language,
  tts_voice,
  media_provider,
  media_video_id,
  media_start_seconds,
  media_end_seconds,
  media_duration_seconds,
  played_at,
  meta
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;