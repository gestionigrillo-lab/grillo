-- Migration: Create audit_logs table
-- Secure Samsung Vault - Audit Trail

CREATE TABLE IF NOT EXISTS audit_logs (
    id         VARCHAR(36)  PRIMARY KEY,
    device_id  VARCHAR(36)  NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    action     VARCHAR(64)  NOT NULL,
    detail     TEXT         DEFAULT '',
    ip         VARCHAR(45)  DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_device ON audit_logs(device_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
