package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/gestionigrillo/secure-samsung-vault/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DeviceService struct {
	devices *repository.DeviceRepository
	audit   *repository.AuditRepository
	redis   *redis.Client
}

func NewDeviceService(devices *repository.DeviceRepository, audit *repository.AuditRepository, redis *redis.Client) *DeviceService {
	return &DeviceService{devices: devices, audit: audit, redis: redis}
}

func (s *DeviceService) Enroll(ctx context.Context, req *models.EnrollRequest, ip string) (*models.EnrollResponse, error) {
	deviceID := uuid.New().String()
	token, err := generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	now := time.Now().UTC()
	device := &models.Device{
		ID:             deviceID,
		DeviceModel:    req.DeviceModel,
		AndroidVersion: req.AndroidVersion,
		AppVersion:     req.AppVersion,
		PushToken:      req.PushToken,
		Status:         models.DeviceStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.devices.Create(ctx, device); err != nil {
		return nil, err
	}

	// Store token in Redis with device mapping
	if err := s.redis.Set(ctx, "device_token:"+token, deviceID, 0).Err(); err != nil {
		return nil, fmt.Errorf("store device token: %w", err)
	}

	s.logAudit(ctx, deviceID, "enroll", "device enrolled", ip)

	return &models.EnrollResponse{DeviceID: deviceID, Token: token}, nil
}

func (s *DeviceService) Heartbeat(ctx context.Context, req *models.HeartbeatRequest, ip string) (*models.HeartbeatResponse, error) {
	device, err := s.devices.GetByID(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	if device.Status == models.DeviceStatusRevoked {
		return &models.HeartbeatResponse{Status: models.DeviceStatusRevoked}, nil
	}

	if err := s.devices.UpdateHeartbeat(ctx, req.DeviceID, req.AppVersion); err != nil {
		return nil, err
	}

	// Log compliance warning if rooted
	if req.IsRooted {
		s.logAudit(ctx, req.DeviceID, "compliance_warning", "device is rooted", ip)
	}

	s.logAudit(ctx, req.DeviceID, "heartbeat", "", ip)

	return &models.HeartbeatResponse{Status: device.Status}, nil
}

func (s *DeviceService) ValidateToken(token string) (string, error) {
	deviceID, err := s.redis.Get(context.Background(), "device_token:"+token).Result()
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	return deviceID, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	return s.devices.GetByID(ctx, id)
}

func (s *DeviceService) UpdateStatus(ctx context.Context, deviceID string, status models.DeviceStatus, ip string) error {
	if err := s.devices.UpdateStatus(ctx, deviceID, status); err != nil {
		return err
	}
	s.logAudit(ctx, deviceID, "status_change", string(status), ip)
	return nil
}

func (s *DeviceService) logAudit(ctx context.Context, deviceID, action, detail, ip string) {
	_ = s.audit.Log(ctx, &models.AuditLog{
		ID:        uuid.New().String(),
		DeviceID:  deviceID,
		Action:    action,
		Detail:    detail,
		IP:        ip,
		CreatedAt: time.Now().UTC(),
	})
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
