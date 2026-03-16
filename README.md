# Secure Samsung Vault - Go Backend

A secure device management backend for Samsung Knox Vault operations. Provides device enrollment, heartbeat monitoring, remote command execution (lock, wipe, revoke), and audit logging.

## Quick Start

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

## API Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | `/health` | Health check | None |
| POST | `/api/v1/device/enroll` | Enroll a device | X-Device-Token |
| POST | `/api/v1/device/heartbeat` | Device heartbeat | X-Device-Token |
| GET | `/api/v1/device/{deviceId}/commands` | Get pending commands | X-Device-Token |
| POST | `/api/v1/device/commands/{commandId}/ack` | Acknowledge command | X-Device-Token |
| POST | `/api/v1/admin/commands` | Issue command to device | Bearer token |

## Configuration

Environment variables (see `.env.example`):

- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `ADMIN_TOKEN` - Admin API authentication token
- `SERVER_PORT` - Server port (default: 8080)
- `JWT_SECRET` - JWT signing secret

## Development

```bash
# Run dependencies
docker-compose up postgres redis -d

# Run the app
export DATABASE_URL=postgres://vault:vault_secret@localhost:5432/samsung_vault?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export ADMIN_TOKEN=dev-token
go run ./cmd/server
```
