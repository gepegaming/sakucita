CREATE TYPE ledger_entry_type AS ENUM (
  'DEPOSIT',
  'WITHDRAW',
  'TRANSFER',
  'FEE',
  'REFUND',
  'ADJUSTMENT'
);

CREATE TYPE ledger_source_type AS ENUM (
  'TRANSACTION',
  'LEDGER'
);

CREATE TABLE wallet_ledger_entries (
  id UUID PRIMARY KEY,

  wallet_id UUID NOT NULL
    REFERENCES wallets(id),

  entry_type ledger_entry_type NOT NULL,

  -- positive / negative allowed
  amount BIGINT NOT NULL CHECK (amount <> 0),
  currency CHAR(3) NOT NULL DEFAULT 'IDR',

  -- traceability
  source_type ledger_source_type NOT NULL,
  source_id UUID NOT NULL,

  -- copied snapshot for audit
  fee_fixed BIGINT NOT NULL DEFAULT 0,
  fee_percentage NUMERIC(5,4) NOT NULL DEFAULT 0,
  fee_amount BIGINT NOT NULL DEFAULT 0,
  net_amount BIGINT NOT NULL,

  description TEXT,
  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_wallet_ledger_wallet
  ON wallet_ledger_entries(wallet_id);

CREATE INDEX idx_wallet_ledger_source
  ON wallet_ledger_entries(source_type, source_id);
