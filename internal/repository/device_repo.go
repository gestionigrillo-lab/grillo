package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/secure-samsung-vault/grillo/internal/models"
)

type DeviceRepository struct {
	pool *pgxpool.Pool
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{pool: pool}
}

func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	query := `
		INSERT INTO devices (id, device_id, model, os_version, public_key, enrollment_token, status, vault_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		device.ID,
		device.DeviceID,
		device.Model,
		device.OSVersion,
		device.PublicKey,
		device.EnrollmentToken,
		device.Status,
		device.VaultStatus,
		device.CreatedAt,
		device.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	return nil
}

func (r *DeviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.Device, error) {
	query := `
		SELECT id, device_id, model, os_version, public_key, enrollment_token, status,
		       last_heartbeat, battery_level, is_locked, vault_status, created_at, updated_at
		FROM devices WHERE device_id = $1`

	var d models.Device
	err := r.pool.QueryRow(ctx, query, deviceID).Scan(
		&d.ID, &d.DeviceID, &d.Model, &d.OSVersion, &d.PublicKey,
		&d.EnrollmentToken, &d.Status, &d.LastHeartbeat, &d.BatteryLevel,
		&d.IsLocked, &d.VaultStatus, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get device by device_id: %w", err)
	}
	return &d, nil
}

func (r *DeviceRepository) GetByEnrollmentToken(ctx context.Context, token string) (*models.Device, error) {
	query := `
		SELECT id, device_id, model, os_version, public_key, enrollment_token, status,
		       last_heartbeat, battery_level, is_locked, vault_status, created_at, updated_at
		FROM devices WHERE enrollment_token = $1`

	var d models.Device
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&d.ID, &d.DeviceID, &d.Model, &d.OSVersion, &d.PublicKey,
		&d.EnrollmentToken, &d.Status, &d.LastHeartbeat, &d.BatteryLevel,
		&d.IsLocked, &d.VaultStatus, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get device by enrollment_token: %w", err)
	}
	return &d, nil
}

func (r *DeviceRepository) UpdateHeartbeat(ctx context.Context, deviceID string, batteryLevel int, isLocked bool, vaultStatus string) error {
	query := `
		UPDATE devices
		SET last_heartbeat = $1, battery_level = $2, is_locked = $3, vault_status = $4, updated_at = $5
		WHERE device_id = $6`

	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, query, now, batteryLevel, isLocked, vaultStatus, now, deviceID)
	if err != nil {
		return fmt.Errorf("update heartbeat: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("device not found: %s", deviceID)
	}
	return nil
}

func (r *DeviceRepository) UpdateStatus(ctx context.Context, deviceID string, status models.DeviceStatus) error {
	query := `UPDATE devices SET status = $1, updated_at = $2 WHERE device_id = $3`
	now := time.Now().UTC()
	result, err := r.pool.Exec(ctx, query, status, now, deviceID)
	if err != nil {
		return fmt.Errorf("update device status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("device not found: %s", deviceID)
	}
	return nil
}

func (r *DeviceRepository) Exists(ctx context.Context, deviceID string) (bool, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, "SELECT id FROM devices WHERE device_id = $1", deviceID).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("check device exists: %w", err)
	}
	return true, nil
}
