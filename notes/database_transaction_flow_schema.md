oke, makasih sudah ngebalikin ke base yang benar. kamu bener: wallet + ledger itu core, donation message cuma layer experience.
sekarang aku rewrite TOTAL jadi FINAL MIGRATION SQL (copy-paste ready), rapih, konsisten, tidak melenceng, dan sesuai flow wallet yang kita sepakati.

Aku susun urut migration yang benar dan tidak ada domain bocor.

✅ URUTAN MIGRATION (WAJIB INI URUTANNYA)
001_enums.sql
002_transactions.sql
003_wallets.sql
004_wallet_ledger_entries.sql
005_donation_messages.sql

📁 001_enums.sql
CREATE TYPE transaction_status AS ENUM (
  'PENDING',
  'PAID',
  'SETTLED',
  'FAILED',
  'REFUNDED'
);

CREATE TYPE wallet_type AS ENUM (
  'CASH',
  'PENDING',
  'PLATFORM'
);

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

CREATE TYPE media_provider AS ENUM (
  'YOUTUBE'
);

CREATE TYPE donation_media_type AS ENUM (
  'NONE',
  'TTS',
  'YOUTUBE'
);

📁 002_transactions.sql

(PURE FINANCIAL SOURCE — TIDAK ADA UX)

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- actors
  payer_user_id UUID NULL,
  payee_user_id UUID NOT NULL,

  -- money (smallest unit)
  amount BIGINT NOT NULL CHECK (amount > 0),

  -- fee config (parameter snapshot)
  fee_fixed BIGINT NOT NULL DEFAULT 0,
  fee_percentage NUMERIC(5,4) NOT NULL DEFAULT 0,

  -- final result
  fee_amount BIGINT NOT NULL,
  net_amount BIGINT NOT NULL,

  -- rounding transparency
  rounding_amount NUMERIC(10,6) NOT NULL DEFAULT 0,
  rounding_strategy VARCHAR(30) NOT NULL DEFAULT 'FLOOR_TO_PLATFORM',

  currency CHAR(3) NOT NULL DEFAULT 'IDR',

  status transaction_status NOT NULL DEFAULT 'PENDING',

  -- external payment reference
  external_reference VARCHAR(100),

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  paid_at TIMESTAMPTZ,
  settled_at TIMESTAMPTZ,

  CHECK (fee_amount + net_amount = amount)
);

CREATE INDEX idx_transactions_payee
  ON transactions(payee_user_id);

CREATE INDEX idx_transactions_status
  ON transactions(status);

CREATE INDEX idx_transactions_external_reference
  ON transactions(external_reference);


📌 Catatan penting

transactions tidak tahu donation

aman untuk audit & ledger replay

refund / chargeback tetap masuk sini

📁 003_wallets.sql

(READ MODEL / PROJECTION)

CREATE TABLE wallets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

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

📁 004_wallet_ledger_entries.sql

(FINANCIAL SOURCE OF TRUTH)

CREATE TABLE wallet_ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

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


📌 Ini ledger beneran

immutable

replayable

audit-safe

wallet balance selalu bisa direbuild

📁 005_donation_messages.sql

(EXPERIENCE LAYER — TERPISAH TOTAL)

CREATE TABLE donation_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  transaction_id UUID NOT NULL
    REFERENCES transactions(id) ON DELETE CASCADE,

  streamer_user_id UUID NOT NULL,

  donor_name VARCHAR(100),
  message TEXT,
  is_anonymous BOOLEAN NOT NULL DEFAULT false,

  media_type donation_media_type NOT NULL DEFAULT 'NONE',

  -- TTS
  tts_language VARCHAR(10),
  tts_voice VARCHAR(50),

  -- Media (YouTube only for now)
  media_provider media_provider,
  media_video_id VARCHAR(50),
  media_start_seconds INT CHECK (media_start_seconds >= 0),
  media_end_seconds INT CHECK (media_end_seconds >= 0),
  media_duration_seconds INT CHECK (media_duration_seconds >= 0),

  -- FACT ONLY (no control state)
  played_at TIMESTAMPTZ,

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (transaction_id)
);

CREATE INDEX idx_donation_messages_streamer_unplayed
  ON donation_messages(streamer_user_id, created_at)
  WHERE played_at IS NULL;


📌 Kenapa ini bersih

donation = turunan transaksi

overlay cuma render

backend controller + Redis queue

DB hanya simpan fakta yang sudah terjadi

🧠 RELATION MODEL (FINAL)
transactions 1 ─── 0..1 donation_messages
transactions 1 ─── N wallet_ledger_entries
wallets      1 ─── N wallet_ledger_entries

🔥 FLOW FINAL (TIDAK RANCU)
Payment Provider
   ↓
transactions (PAID)
   ↓
wallet_ledger_entries (atomic tx)
   ↓
wallets.balance (projection)
   ↓
donation_messages (UX fact)
   ↓
Redis queue
   ↓
Overlay render
   ↓
WS FINISHED
   ↓
UPDATE donation_messages.played_at

✅ CHECKLIST “MENTOR MODE”

✔ wallet & ledger tetap utuh
✔ donation domain bersih
✔ overlay bodoh
✔ backend authoritative
✔ redis hanya runtime
✔ monolith → microservice ready
✔ audit & replay aman

kalau mau lanjut, step paling tepat berikutnya menurutku:

SQL rebuild wallet balance

Redis queue schema (ZSET)

Donation dispatcher state machine

Refund + ledger reversal

sqlc + Go repo layer

tinggal bilang, kita gas satu-satu 🔥