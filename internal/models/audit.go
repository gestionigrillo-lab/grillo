package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID       `json:"id"`
	DeviceID  string          `json:"device_id,omitempty"`
	EventType string          `json:"event_type"`
	Details   json.RawMessage `json:"details,omitempty"`
	IPAddress string          `json:"ip_address,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}
