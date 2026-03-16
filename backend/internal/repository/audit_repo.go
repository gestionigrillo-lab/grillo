package repository

import (
	"context"
	"fmt"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Log(ctx context.Context, entry *models.AuditLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_logs (id, device_id, action, detail, ip, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		entry.ID, entry.DeviceID, entry.Action, entry.Detail, entry.IP, entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) GetByDevice(ctx context.Context, deviceID string, limit int) ([]models.AuditLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, device_id, action, detail, ip, created_at
		 FROM audit_logs WHERE device_id = $1 ORDER BY created_at DESC LIMIT $2`,
		deviceID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.DeviceID, &l.Action, &l.Detail, &l.IP, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}
