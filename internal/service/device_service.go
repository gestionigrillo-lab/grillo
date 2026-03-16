package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/repository"
)

type DeviceService struct {
	deviceRepo *repository.DeviceRepository
	auditSvc   *AuditService
}

func NewDeviceService(deviceRepo *repository.DeviceRepository, auditSvc *AuditService) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
		auditSvc:   auditSvc,
	}
}

func (s *DeviceService) Enroll(ctx context.Context, req models.EnrollRequest, ipAddress string) (*models.EnrollResponse, error) {
	if req.DeviceID == "" {
		return nil, fmt.Errorf("deviceId is required")
	}

	exists, err := s.deviceRepo.Exists(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("check device existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("device already enrolled")
	}

	token, err := generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate enrollment token: %w", err)
	}

	now := time.Now().UTC()
	device := &models.Device{
		ID:              uuid.New(),
		DeviceID:        req.DeviceID,
		Model:           req.Model,
		OSVersion:       req.OSVersion,
		PublicKey:        req.PublicKey,
		EnrollmentToken: token,
		Status:          models.DeviceStatusActive,
		VaultStatus:     "sealed",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	s.auditSvc.Log(ctx, req.DeviceID, "device_enrolled", map[string]string{
		"model":      req.Model,
		"os_version": req.OSVersion,
	}, ipAddress)

	return &models.EnrollResponse{
		EnrollmentToken: token,
		DeviceID:        req.DeviceID,
	}, nil
}

func (s *DeviceService) Heartbeat(ctx context.Context, req models.HeartbeatRequest, ipAddress string) (*models.HeartbeatResponse, error) {
	if req.DeviceID == "" {
		return nil, fmt.Errorf("deviceId is required")
	}

	err := s.deviceRepo.UpdateHeartbeat(ctx, req.DeviceID, req.BatteryLevel, req.IsLocked, req.VaultStatus)
	if err != nil {
		return nil, fmt.Errorf("update heartbeat: %w", err)
	}

	s.auditSvc.Log(ctx, req.DeviceID, "heartbeat", map[string]interface{}{
		"battery_level": req.BatteryLevel,
		"is_locked":     req.IsLocked,
		"vault_status":  req.VaultStatus,
	}, ipAddress)

	return &models.HeartbeatResponse{
		Status:  "ok",
		Message: "heartbeat recorded",
	}, nil
}

func (s *DeviceService) GetByEnrollmentToken(ctx context.Context, token string) (*models.Device, error) {
	return s.deviceRepo.GetByEnrollmentToken(ctx, token)
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
