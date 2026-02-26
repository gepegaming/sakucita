-- transaction status
CREATE TYPE transaction_status AS ENUM (
  'INITIAL',
  'PENDING',
  'PAID',
  'FAILED',
  'EXPIRED',
  'REFUNDED'
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY,

  donation_message_id UUID NOT NULL
    REFERENCES donation_messages(id)
    ON DELETE RESTRICT,

  payment_channel_id INT NOT NULL
    REFERENCES payment_channels(id),

  payer_user_id UUID,
  payee_user_id UUID NOT NULL,

  -- gross donation amount
  gross_paid_amount BIGINT NOT NULL CHECK (gross_paid_amount > 0),

  ------------------------------------------------------------------
  -- GATEWAY SNAPSHOT
  ------------------------------------------------------------------
  gateway_fee_fixed BIGINT NOT NULL DEFAULT 0,
  gateway_fee_percentage_bps INT NOT NULL DEFAULT 0,
  gateway_fee_amount BIGINT NOT NULL DEFAULT 0,

  ------------------------------------------------------------------
  -- PLATFORM SNAPSHOT
  ------------------------------------------------------------------
  platform_fee_fixed BIGINT NOT NULL DEFAULT 0,
  platform_fee_percentage_bps INT NOT NULL DEFAULT 0,
  platform_fee_amount BIGINT NOT NULL DEFAULT 0,

  ------------------------------------------------------------------
  -- TOTAL FEE (gateway + platform)
  ------------------------------------------------------------------
  fee_fixed BIGINT NOT NULL DEFAULT 0,
  fee_percentage_bps INT NOT NULL DEFAULT 0,
  fee_amount BIGINT NOT NULL,

  -- amount diterima creator setelah di potong fee pg dan platform jd adil ye
  net_amount BIGINT NOT NULL,

  currency CHAR(3) NOT NULL DEFAULT 'IDR',

  status transaction_status NOT NULL DEFAULT 'PENDING',

  external_reference VARCHAR(255) UNIQUE,

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  paid_at TIMESTAMPTZ,
  settled_at TIMESTAMPTZ,

  CHECK (fee_amount + net_amount = gross_paid_amount),
  CHECK (gateway_fee_amount + platform_fee_amount = fee_amount),
  UNIQUE (donation_message_id)
);
