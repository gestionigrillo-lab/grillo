-- Migration: Create devices table
-- Secure Samsung Vault - Device Security Layer

CREATE TABLE IF NOT EXISTS devices (
    id              VARCHAR(36) PRIMARY KEY,
    device_model    VARCHAR(128) NOT NULL,
    android_version VARCHAR(32)  NOT NULL,
    app_version     VARCHAR(32)  NOT NULL,
    push_token      TEXT         DEFAULT '',
    status          VARCHAR(16)  NOT NULL DEFAULT 'active',
    last_heartbeat  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_last_heartbeat ON devices(last_heartbeat);
