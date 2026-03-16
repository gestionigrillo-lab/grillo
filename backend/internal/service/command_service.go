package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/gestionigrillo/secure-samsung-vault/internal/repository"
	"github.com/google/uuid"
)

type CommandService struct {
	commands *repository.CommandRepository
	devices  *repository.DeviceRepository
	audit    *repository.AuditRepository
}

func NewCommandService(commands *repository.CommandRepository, devices *repository.DeviceRepository, audit *repository.AuditRepository) *CommandService {
	return &CommandService{commands: commands, devices: devices, audit: audit}
}

func (s *CommandService) CreateCommand(ctx context.Context, req *models.CreateCommandRequest, ip string) (*models.Command, error) {
	// Verify device exists
	if _, err := s.devices.GetByID(ctx, req.DeviceID); err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	cmd := &models.Command{
		ID:        uuid.New().String(),
		DeviceID:  req.DeviceID,
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    models.CommandStatusPending,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.commands.Create(ctx, cmd); err != nil {
		return nil, err
	}

	_ = s.audit.Log(ctx, &models.AuditLog{
		ID:        uuid.New().String(),
		DeviceID:  req.DeviceID,
		Action:    "command_created",
		Detail:    string(req.Type),
		IP:        ip,
		CreatedAt: time.Now().UTC(),
	})

	return cmd, nil
}

func (s *CommandService) GetPendingCommands(ctx context.Context, deviceID string) ([]models.Command, error) {
	return s.commands.GetPendingByDevice(ctx, deviceID)
}

func (s *CommandService) AcknowledgeCommand(ctx context.Context, commandID string, req *models.AckCommandRequest, ip string) error {
	cmd, err := s.commands.GetByID(ctx, commandID)
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	if err := s.commands.Acknowledge(ctx, commandID, req.Success); err != nil {
		return err
	}

	// Apply side effects for successful commands
	if req.Success {
		switch cmd.Type {
		case models.CommandLock:
			_ = s.devices.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusLocked)
		case models.CommandRevoke:
			_ = s.devices.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusRevoked)
		case models.CommandWipe:
			_ = s.devices.UpdateStatus(ctx, cmd.DeviceID, models.DeviceStatusWiped)
		}
	}

	action := "command_acked"
	if !req.Success {
		action = "command_failed"
	}

	_ = s.audit.Log(ctx, &models.AuditLog{
		ID:        uuid.New().String(),
		DeviceID:  cmd.DeviceID,
		Action:    action,
		Detail:    fmt.Sprintf("%s: %s", cmd.Type, req.Message),
		IP:        ip,
		CreatedAt: time.Now().UTC(),
	})

	return nil
}
