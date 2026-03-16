package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/repository"
)

type AuditService struct {
	auditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) Log(ctx context.Context, deviceID string, eventType string, details interface{}, ipAddress string) {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		log.Printf("failed to marshal audit details: %v", err)
		return
	}

	entry := &models.AuditLog{
		ID:        uuid.New(),
		DeviceID:  deviceID,
		EventType: eventType,
		Details:   detailsJSON,
		IPAddress: ipAddress,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.auditRepo.Create(ctx, entry); err != nil {
		log.Printf("failed to create audit log: %v", err)
	}
}

func (s *AuditService) GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.auditRepo.GetByDeviceID(ctx, deviceID, limit)
}
