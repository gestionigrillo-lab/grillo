package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceStatus string

const (
	DeviceStatusActive  DeviceStatus = "active"
	DeviceStatusLocked  DeviceStatus = "locked"
	DeviceStatusRevoked DeviceStatus = "revoked"
	DeviceStatusWiped   DeviceStatus = "wiped"
)

type Device struct {
	ID              uuid.UUID    `json:"id"`
	DeviceID        string       `json:"device_id"`
	Model           string       `json:"model"`
	OSVersion       string       `json:"os_version"`
	PublicKey        string       `json:"public_key"`
	EnrollmentToken string       `json:"enrollment_token"`
	Status          DeviceStatus `json:"status"`
	LastHeartbeat   *time.Time   `json:"last_heartbeat,omitempty"`
	BatteryLevel    int          `json:"battery_level"`
	IsLocked        bool         `json:"is_locked"`
	VaultStatus     string       `json:"vault_status"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type EnrollRequest struct {
	DeviceID  string `json:"deviceId"`
	Model     string `json:"model"`
	OSVersion string `json:"osVersion"`
	PublicKey string `json:"publicKey"`
}

type EnrollResponse struct {
	EnrollmentToken string `json:"enrollmentToken"`
	DeviceID        string `json:"deviceId"`
}

type HeartbeatRequest struct {
	DeviceID     string `json:"deviceId"`
	BatteryLevel int    `json:"batteryLevel"`
	IsLocked     bool   `json:"isLocked"`
	VaultStatus  string `json:"vaultStatus"`
	Timestamp    string `json:"timestamp"`
}

type HeartbeatResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
