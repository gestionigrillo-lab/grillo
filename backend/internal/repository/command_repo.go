package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandRepository struct {
	pool *pgxpool.Pool
}

func NewCommandRepository(pool *pgxpool.Pool) *CommandRepository {
	return &CommandRepository{pool: pool}
}

func (r *CommandRepository) Create(ctx context.Context, cmd *models.Command) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO device_commands (id, device_id, type, payload, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		cmd.ID, cmd.DeviceID, cmd.Type, cmd.Payload, cmd.Status, cmd.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert command: %w", err)
	}
	return nil
}

func (r *CommandRepository) GetPendingByDevice(ctx context.Context, deviceID string) ([]models.Command, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, device_id, type, payload, status, created_at, acked_at
		 FROM device_commands WHERE device_id = $1 AND status = $2 ORDER BY created_at ASC`,
		deviceID, models.CommandStatusPending,
	)
	if err != nil {
		return nil, fmt.Errorf("query pending commands: %w", err)
	}
	defer rows.Close()

	var commands []models.Command
	for rows.Next() {
		var c models.Command
		if err := rows.Scan(&c.ID, &c.DeviceID, &c.Type, &c.Payload, &c.Status, &c.CreatedAt, &c.AckedAt); err != nil {
			return nil, fmt.Errorf("scan command: %w", err)
		}
		commands = append(commands, c)
	}
	return commands, nil
}

func (r *CommandRepository) Acknowledge(ctx context.Context, commandID string, success bool) error {
	now := time.Now().UTC()
	status := models.CommandStatusAcked
	if !success {
		status = models.CommandStatusFailed
	}

	_, err := r.pool.Exec(ctx,
		`UPDATE device_commands SET status = $1, acked_at = $2 WHERE id = $3`,
		status, now, commandID,
	)
	if err != nil {
		return fmt.Errorf("ack command: %w", err)
	}
	return nil
}

func (r *CommandRepository) GetByID(ctx context.Context, id string) (*models.Command, error) {
	c := &models.Command{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, device_id, type, payload, status, created_at, acked_at
		 FROM device_commands WHERE id = $1`, id,
	).Scan(&c.ID, &c.DeviceID, &c.Type, &c.Payload, &c.Status, &c.CreatedAt, &c.AckedAt)
	if err != nil {
		return nil, fmt.Errorf("get command: %w", err)
	}
	return c, nil
}
