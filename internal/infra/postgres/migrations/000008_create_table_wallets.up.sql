CREATE TYPE wallet_type AS ENUM (
  'CASH',
  'PENDING',
  'PLATFORM'
);

CREATE TABLE wallets (
  id UUID PRIMARY KEY,

  user_id UUID NOT NULL,
  type wallet_type NOT NULL,

  name VARCHAR(50) NOT NULL,
  slug VARCHAR(50) NOT NULL,

  -- derived value (rebuildable)
  balance BIGINT NOT NULL DEFAULT 0,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (user_id, type),
  UNIQUE (user_id, slug)
);