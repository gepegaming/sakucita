CREATE TYPE donation_media_type AS ENUM (
  'TEXT',
  'YOUTUBE',
  'GIF'
);

CREATE TABLE donation_messages (
  id UUID PRIMARY KEY,

  payee_user_id UUID NOT NULL,

  payer_user_id UUID,
  payer_name    VARCHAR(100) NOT NULL,
  email         VARCHAR(255),
  message       TEXT,

  media_type donation_media_type NOT NULL,

  -- media (nullable)
  media_url TEXT,              -- youtube url / gif url
  media_start_seconds INT NOT NULL DEFAULT 0,

  -- hasil kalkulasi backend (bukan input user)
  charged_seconds INT,         -- detik yang DICASHT (<=300)
  price_per_second BIGINT,     -- misal 500
  amount BIGINT NOT NULL,      -- charged_seconds * price_per_second
  currency CHAR(3) NOT NULL DEFAULT 'IDR',   -- IDR, dll

  status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
  -- CREATED | PAID | PLAYED | CANCELED | REJECTED

  meta        JSONB,
  played_at   TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
