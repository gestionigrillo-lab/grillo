package models

import "time"

type CommandType string

const (
	CommandLock   CommandType = "lock"
	CommandWipe   CommandType = "wipe"
	CommandRevoke CommandType = "revoke"
)

type CommandStatus string

const (
	CommandStatusPending CommandStatus = "pending"
	CommandStatusAcked   CommandStatus = "acked"
	CommandStatusFailed  CommandStatus = "failed"
)

type Command struct {
	ID        string        `json:"id"`
	DeviceID  string        `json:"device_id"`
	Type      CommandType   `json:"type"`
	Payload   string        `json:"payload,omitempty"`
	Status    CommandStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	AckedAt   *time.Time    `json:"acked_at,omitempty"`
}

type CreateCommandRequest struct {
	DeviceID string      `json:"device_id"`
	Type     CommandType `json:"type"`
	Payload  string      `json:"payload,omitempty"`
}

type AckCommandRequest struct {
	DeviceID string `json:"device_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
}
