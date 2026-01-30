CREATE TYPE media_provider AS ENUM (
  'YOUTUBE',
  'GIF'
);

CREATE TYPE donation_media_type AS ENUM (
  'TEXT',
  'TTS',
  'YOUTUBE',
  'GIF'
);

CREATE TABLE donation_messages (
  id UUID PRIMARY KEY,

  payee_user_id  UUID NOT NULL,

  payer_user_id     UUID,
  payer_name        VARCHAR(100) NOT NULL,
  message           TEXT NOT NULL,
  email             VARCHAR(100) NOT NULL,


  media_type donation_media_type NOT NULL DEFAULT 'TEXT',

  -- TTS
  tts_language VARCHAR(10),
  tts_voice VARCHAR(50),

  -- Media
  media_provider media_provider,
  media_video_id VARCHAR(50),
  media_start_seconds INT CHECK (media_start_seconds >= 0),
  media_end_seconds INT CHECK (media_end_seconds >= 0),
  media_duration_seconds INT CHECK (media_duration_seconds >= 0),

  -- playback fact
  played_at TIMESTAMPTZ,

  -- lifecycle ringan (opsional tapi berguna)
  status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
  -- CREATED | REJECTED | WAITING_PAYMENT | READY | PLAYED

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

