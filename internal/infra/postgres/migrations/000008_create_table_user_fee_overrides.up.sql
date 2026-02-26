CREATE TABLE user_fee_overrides (
  id UUID PRIMARY KEY,

  user_id UUID NOT NULL REFERENCES users(id),
  payment_channel_id INT NOT NULL REFERENCES payment_channels(id),

  platform_fee_fixed BIGINT NOT NULL DEFAULT 0,
  platform_fee_percentage_bps INT NOT NULL DEFAULT 0,

  UNIQUE (user_id, payment_channel_id)
);
