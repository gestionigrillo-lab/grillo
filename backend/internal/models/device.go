package models

import "time"

type DeviceStatus string

const (
	DeviceStatusActive  DeviceStatus = "active"
	DeviceStatusLocked  DeviceStatus = "locked"
	DeviceStatusRevoked DeviceStatus = "revoked"
	DeviceStatusWiped   DeviceStatus = "wiped"
)

type Device struct {
	ID            string       `json:"id"`
	DeviceModel   string       `json:"device_model"`
	AndroidVersion string      `json:"android_version"`
	AppVersion    string       `json:"app_version"`
	PushToken     string       `json:"push_token,omitempty"`
	Status        DeviceStatus `json:"status"`
	LastHeartbeat *time.Time   `json:"last_heartbeat,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type EnrollRequest struct {
	DeviceModel    string `json:"device_model"`
	AndroidVersion string `json:"android_version"`
	AppVersion     string `json:"app_version"`
	PushToken      string `json:"push_token,omitempty"`
}

type EnrollResponse struct {
	DeviceID string `json:"device_id"`
	Token    string `json:"token"`
}

type HeartbeatRequest struct {
	DeviceID   string `json:"device_id"`
	AppVersion string `json:"app_version"`
	BatteryPct int    `json:"battery_pct,omitempty"`
	IsRooted   bool   `json:"is_rooted"`
}

type HeartbeatResponse struct {
	Status   DeviceStatus `json:"status"`
	Commands []Command    `json:"pending_commands,omitempty"`
}
