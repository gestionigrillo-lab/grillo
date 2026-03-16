CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    device_id TEXT UNIQUE NOT NULL,
    model TEXT,
    os_version TEXT,
    public_key TEXT,
    enrollment_token TEXT UNIQUE,
    status TEXT DEFAULT 'active',
    last_heartbeat TIMESTAMPTZ,
    battery_level INT DEFAULT 0,
    is_locked BOOLEAN DEFAULT false,
    vault_status TEXT DEFAULT 'sealed',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_devices_device_id ON devices(device_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_enrollment_token ON devices(enrollment_token);
