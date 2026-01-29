CREATE TABLE auth_identities (
  id             UUID PRIMARY KEY,
  user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  provider       VARCHAR(20) NOT NULL,   -- local, google, phone
  provider_id    VARCHAR(255) NOT NULL,  -- email / google sub / phone

  password_hash  VARCHAR(255),           -- hanya local
  totp_secret    VARCHAR(255),
  totp_enabled   BOOLEAN NOT NULL DEFAULT FALSE,

  meta           JSONB,
  last_login_at  TIMESTAMPTZ,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ,

  CONSTRAINT uq_provider_providerid UNIQUE (provider, provider_id),
  CONSTRAINT uq_user_provider UNIQUE (user_id, provider)
);

CREATE INDEX idx_auth_identities_user_provider
  ON auth_identities(user_id, provider);