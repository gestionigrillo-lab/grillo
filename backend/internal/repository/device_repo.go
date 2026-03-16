package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository struct {
	pool *pgxpool.Pool
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{pool: pool}
}

func (r *DeviceRepository) Create(ctx context.Context, d *models.Device) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO devices (id, device_model, android_version, app_version, push_token, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		d.ID, d.DeviceModel, d.AndroidVersion, d.AppVersion, d.PushToken, d.Status, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	return nil
}

func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	d := &models.Device{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, device_model, android_version, app_version, push_token, status, last_heartbeat, created_at, updated_at
		 FROM devices WHERE id = $1`, id,
	).Scan(&d.ID, &d.DeviceModel, &d.AndroidVersion, &d.AppVersion, &d.PushToken, &d.Status, &d.LastHeartbeat, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	return d, nil
}

func (r *DeviceRepository) UpdateHeartbeat(ctx context.Context, id string, appVersion string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET last_heartbeat = $1, app_version = $2, updated_at = $3 WHERE id = $4`,
		now, appVersion, now, id,
	)
	if err != nil {
		return fmt.Errorf("update heartbeat: %w", err)
	}
	return nil
}

func (r *DeviceRepository) UpdateStatus(ctx context.Context, id string, status models.DeviceStatus) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET status = $1, updated_at = $2 WHERE id = $3`,
		status, now, id,
	)
	if err != nil {
		return fmt.Errorf("update device status: %w", err)
	}
	return nil
}
