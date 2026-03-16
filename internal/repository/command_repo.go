package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/secure-samsung-vault/grillo/internal/models"
)

type CommandRepository struct {
	pool *pgxpool.Pool
}

func NewCommandRepository(pool *pgxpool.Pool) *CommandRepository {
	return &CommandRepository{pool: pool}
}

func (r *CommandRepository) Create(ctx context.Context, cmd *models.DeviceCommand) error {
	query := `
		INSERT INTO device_commands (id, device_id, command_type, payload, status, issued_by, issued_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		cmd.ID, cmd.DeviceID, cmd.CommandType, cmd.Payload,
		cmd.Status, cmd.IssuedBy, cmd.IssuedAt,
	)
	if err != nil {
		return fmt.Errorf("insert command: %w", err)
	}
	return nil
}

func (r *CommandRepository) GetPendingByDeviceID(ctx context.Context, deviceID string) ([]models.DeviceCommand, error) {
	query := `
		SELECT id, device_id, command_type, payload, status, issued_by, issued_at,
		       delivered_at, executed_at, ack_status, ack_message
		FROM device_commands
		WHERE device_id = $1 AND status = 'pending'
		ORDER BY issued_at ASC`

	rows, err := r.pool.Query(ctx, query, deviceID)
	if err != nil {
		return nil, fmt.Errorf("get pending commands: %w", err)
	}
	defer rows.Close()

	var commands []models.DeviceCommand
	for rows.Next() {
		var cmd models.DeviceCommand
		var issuedBy, ackStatus, ackMessage *string
		err := rows.Scan(
			&cmd.ID, &cmd.DeviceID, &cmd.CommandType, &cmd.Payload,
			&cmd.Status, &issuedBy, &cmd.IssuedAt,
			&cmd.DeliveredAt, &cmd.ExecutedAt, &ackStatus, &ackMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("scan command row: %w", err)
		}
		if issuedBy != nil {
			cmd.IssuedBy = *issuedBy
		}
		if ackStatus != nil {
			cmd.AckStatus = *ackStatus
		}
		if ackMessage != nil {
			cmd.AckMessage = *ackMessage
		}
		commands = append(commands, cmd)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate command rows: %w", err)
	}

	return commands, nil
}

func (r *CommandRepository) GetByID(ctx context.Context, commandID uuid.UUID) (*models.DeviceCommand, error) {
	query := `
		SELECT id, device_id, command_type, payload, status, issued_by, issued_at,
		       delivered_at, executed_at, ack_status, ack_message
		FROM device_commands WHERE id = $1`

	var cmd models.DeviceCommand
	var issuedBy, ackStatus, ackMessage *string
	err := r.pool.QueryRow(ctx, query, commandID).Scan(
		&cmd.ID, &cmd.DeviceID, &cmd.CommandType, &cmd.Payload,
		&cmd.Status, &issuedBy, &cmd.IssuedAt,
		&cmd.DeliveredAt, &cmd.ExecutedAt, &ackStatus, &ackMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("get command by id: %w", err)
	}
	if issuedBy != nil {
		cmd.IssuedBy = *issuedBy
	}
	if ackStatus != nil {
		cmd.AckStatus = *ackStatus
	}
	if ackMessage != nil {
		cmd.AckMessage = *ackMessage
	}
	return &cmd, nil
}

func (r *CommandRepository) MarkDelivered(ctx context.Context, commandID uuid.UUID) error {
	query := `UPDATE device_commands SET status = 'delivered', delivered_at = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, time.Now().UTC(), commandID)
	if err != nil {
		return fmt.Errorf("mark command delivered: %w", err)
	}
	return nil
}

func (r *CommandRepository) Acknowledge(ctx context.Context, commandID uuid.UUID, ackStatus string, executedAt time.Time) error {
	query := `
		UPDATE device_commands
		SET status = 'executed', executed_at = $1, ack_status = $2, ack_message = 'acknowledged'
		WHERE id = $3`

	result, err := r.pool.Exec(ctx, query, executedAt, ackStatus, commandID)
	if err != nil {
		return fmt.Errorf("acknowledge command: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("command not found: %s", commandID)
	}
	return nil
}
