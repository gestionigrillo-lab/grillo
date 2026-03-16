package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CommandType string

const (
	CommandTypeLock   CommandType = "lock"
	CommandTypeWipe   CommandType = "wipe"
	CommandTypeRevoke CommandType = "revoke"
	CommandTypeUnlock CommandType = "unlock"
)

type CommandStatus string

const (
	CommandStatusPending   CommandStatus = "pending"
	CommandStatusDelivered CommandStatus = "delivered"
	CommandStatusExecuted  CommandStatus = "executed"
	CommandStatusFailed    CommandStatus = "failed"
)

type DeviceCommand struct {
	ID          uuid.UUID       `json:"id"`
	DeviceID    string          `json:"device_id"`
	CommandType CommandType     `json:"command_type"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	Status      CommandStatus   `json:"status"`
	IssuedBy    string          `json:"issued_by,omitempty"`
	IssuedAt    time.Time       `json:"issued_at"`
	DeliveredAt *time.Time      `json:"delivered_at,omitempty"`
	ExecutedAt  *time.Time      `json:"executed_at,omitempty"`
	AckStatus   string          `json:"ack_status,omitempty"`
	AckMessage  string          `json:"ack_message,omitempty"`
}

type IssueCommandRequest struct {
	DeviceID    string          `json:"deviceId"`
	CommandType CommandType     `json:"commandType"`
	Payload     json.RawMessage `json:"payload,omitempty"`
}

type IssueCommandResponse struct {
	CommandID   uuid.UUID   `json:"commandId"`
	DeviceID    string      `json:"deviceId"`
	CommandType CommandType `json:"commandType"`
	Status      string      `json:"status"`
}

type AckCommandRequest struct {
	Status     string `json:"status"`
	ExecutedAt string `json:"executedAt"`
}

type AckCommandResponse struct {
	CommandID string `json:"commandId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
