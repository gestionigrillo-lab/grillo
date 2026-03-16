package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/repository"
)

type CommandService struct {
	commandRepo *repository.CommandRepository
	deviceRepo  *repository.DeviceRepository
	auditSvc    *AuditService
}

func NewCommandService(commandRepo *repository.CommandRepository, deviceRepo *repository.DeviceRepository, auditSvc *AuditService) *CommandService {
	return &CommandService{
		commandRepo: commandRepo,
		deviceRepo:  deviceRepo,
		auditSvc:    auditSvc,
	}
}

func (s *CommandService) IssueCommand(ctx context.Context, req models.IssueCommandRequest, issuedBy string, ipAddress string) (*models.IssueCommandResponse, error) {
	if req.DeviceID == "" {
		return nil, fmt.Errorf("deviceId is required")
	}

	validTypes := map[models.CommandType]bool{
		models.CommandTypeLock:   true,
		models.CommandTypeWipe:   true,
		models.CommandTypeRevoke: true,
		models.CommandTypeUnlock: true,
	}
	if !validTypes[req.CommandType] {
		return nil, fmt.Errorf("invalid command type: %s", req.CommandType)
	}

	exists, err := s.deviceRepo.Exists(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("check device: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("device not found: %s", req.DeviceID)
	}

	cmd := &models.DeviceCommand{
		ID:          uuid.New(),
		DeviceID:    req.DeviceID,
		CommandType: req.CommandType,
		Payload:     req.Payload,
		Status:      models.CommandStatusPending,
		IssuedBy:    issuedBy,
		IssuedAt:    time.Now().UTC(),
	}

	if err := s.commandRepo.Create(ctx, cmd); err != nil {
		return nil, fmt.Errorf("create command: %w", err)
	}

	s.auditSvc.Log(ctx, req.DeviceID, "command_issued", map[string]string{
		"command_type": string(req.CommandType),
		"command_id":   cmd.ID.String(),
		"issued_by":    issuedBy,
	}, ipAddress)

	return &models.IssueCommandResponse{
		CommandID:   cmd.ID,
		DeviceID:    req.DeviceID,
		CommandType: req.CommandType,
		Status:      string(models.CommandStatusPending),
	}, nil
}

func (s *CommandService) GetPendingCommands(ctx context.Context, deviceID string) ([]models.DeviceCommand, error) {
	commands, err := s.commandRepo.GetPendingByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("get pending commands: %w", err)
	}

	// Mark commands as delivered
	for _, cmd := range commands {
		_ = s.commandRepo.MarkDelivered(ctx, cmd.ID)
	}

	return commands, nil
}

func (s *CommandService) AcknowledgeCommand(ctx context.Context, commandID uuid.UUID, req models.AckCommandRequest, ipAddress string) (*models.AckCommandResponse, error) {
	cmd, err := s.commandRepo.GetByID(ctx, commandID)
	if err != nil {
		return nil, fmt.Errorf("get command: %w", err)
	}

	executedAt := time.Now().UTC()
	if req.ExecutedAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.ExecutedAt)
		if err == nil {
			executedAt = parsed
		}
	}

	if err := s.commandRepo.Acknowledge(ctx, commandID, req.Status, executedAt); err != nil {
		return nil, fmt.Errorf("acknowledge command: %w", err)
	}

	// Update device status based on command type
	switch cmd.CommandType {
	case models.CommandTypeLock:
		_ = s.deviceRepo.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusLocked)
	case models.CommandTypeWipe:
		_ = s.deviceRepo.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusWiped)
	case models.CommandTypeRevoke:
		_ = s.deviceRepo.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusRevoked)
	case models.CommandTypeUnlock:
		_ = s.deviceRepo.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusActive)
	}

	s.auditSvc.Log(ctx, cmd.DeviceID, "command_acknowledged", map[string]string{
		"command_id":   commandID.String(),
		"command_type": string(cmd.CommandType),
		"ack_status":   req.Status,
	}, ipAddress)

	return &models.AckCommandResponse{
		CommandID: commandID.String(),
		Status:    "acknowledged",
		Message:   "command execution acknowledged",
	}, nil
}
