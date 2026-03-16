-- Migration: Create device_commands table
-- Secure Samsung Vault - Remote Commands

CREATE TABLE IF NOT EXISTS device_commands (
    id         VARCHAR(36)  PRIMARY KEY,
    device_id  VARCHAR(36)  NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    type       VARCHAR(32)  NOT NULL,
    payload    TEXT         DEFAULT '',
    status     VARCHAR(16)  NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    acked_at   TIMESTAMPTZ
);

CREATE INDEX idx_device_commands_device_status ON device_commands(device_id, status);
CREATE INDEX idx_device_commands_created ON device_commands(created_at);
