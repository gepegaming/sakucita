CREATE TABLE users (
  id              UUID PRIMARY KEY,

  email           VARCHAR(255) UNIQUE NOT NULL,
  email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
  phone           VARCHAR(20) UNIQUE,

  name            VARCHAR(100) NOT NULL,
  nickname        VARCHAR(50) UNIQUE NOT NULL,
  image_url       VARCHAR(255),

  single_session  BOOLEAN NOT NULL DEFAULT FALSE,
  
  meta            JSONB,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_email_verified ON users(email_verified);