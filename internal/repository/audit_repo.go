package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/secure-samsung-vault/grillo/internal/models"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, device_id, event_type, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		log.ID, log.DeviceID, log.EventType, log.Details, log.IPAddress, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.AuditLog, error) {
	query := `
		SELECT id, device_id, event_type, details, ip_address, created_at
		FROM audit_logs WHERE device_id = $1
		ORDER BY created_at DESC LIMIT $2`

	rows, err := r.pool.Query(ctx, query, deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("get audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		var deviceIDPtr, ipAddr *string
		err := rows.Scan(&l.ID, &deviceIDPtr, &l.EventType, &l.Details, &ipAddr, &l.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		if deviceIDPtr != nil {
			l.DeviceID = *deviceIDPtr
		}
		if ipAddr != nil {
			l.IPAddress = *ipAddr
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
