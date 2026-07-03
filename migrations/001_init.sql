CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS api_keys (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(100) NOT NULL,
    key_hash     VARCHAR(64)  NOT NULL UNIQUE,
    scopes       TEXT[]       NOT NULL DEFAULT '{"mail:send"}',
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS templates (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL UNIQUE,
    subject    TEXT         NOT NULL,
    html_body  TEXT         NOT NULL,
    text_body  TEXT,
    variables  TEXT[]       NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_logs (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_key_id    UUID        REFERENCES api_keys(id) ON DELETE SET NULL,
    from_address  TEXT        NOT NULL,
    to_addresses  TEXT[]      NOT NULL,
    cc_addresses  TEXT[],
    bcc_addresses TEXT[],
    subject       TEXT        NOT NULL,
    template_name TEXT,
    status        VARCHAR(20) NOT NULL DEFAULT 'queued',
    error_msg     TEXT,
    attempts      INT         NOT NULL DEFAULT 0,
    sent_at       TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_logs_status     ON email_logs(status);
CREATE INDEX IF NOT EXISTS idx_email_logs_api_key    ON email_logs(api_key_id);
CREATE INDEX IF NOT EXISTS idx_email_logs_created_at ON email_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash         ON api_keys(key_hash);
