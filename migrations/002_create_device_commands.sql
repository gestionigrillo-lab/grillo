CREATE TABLE IF NOT EXISTS device_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id TEXT NOT NULL REFERENCES devices(device_id),
    command_type TEXT NOT NULL,
    payload JSONB,
    status TEXT DEFAULT 'pending',
    issued_by TEXT,
    issued_at TIMESTAMPTZ DEFAULT NOW(),
    delivered_at TIMESTAMPTZ,
    executed_at TIMESTAMPTZ,
    ack_status TEXT,
    ack_message TEXT
);

CREATE INDEX idx_device_commands_device_id ON device_commands(device_id);
CREATE INDEX idx_device_commands_status ON device_commands(status);
CREATE INDEX idx_device_commands_device_pending ON device_commands(device_id, status) WHERE status = 'pending';
