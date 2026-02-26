-- name: GetUserFee :one
SELECT
    pc.id AS payment_channel_id,

    COALESCE(ufo.platform_fee_fixed, pc.platform_fee_fixed) 
        AS platform_fee_fixed,

    COALESCE(ufo.platform_fee_percentage_bps, pc.platform_fee_percentage_bps) 
        AS platform_fee_percentage_bps,

    pc.gateway_fee_fixed,
    pc.gateway_fee_percentage_bps

FROM payment_channels pc
LEFT JOIN user_fee_overrides ufo
    ON ufo.payment_channel_id = pc.id
    AND ufo.user_id = sqlc.arg(UserID)

WHERE pc.id = sqlc.arg(PaymentChannelID)
LIMIT 1;

-- name: GetAllUserFeeOverridesByUserID :many
SELECT * FROM user_fee_overrides 
WHERE user_id = $1;