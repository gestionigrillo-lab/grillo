package models

import "time"

type AuditLog struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail,omitempty"`
	IP        string    `json:"ip,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
