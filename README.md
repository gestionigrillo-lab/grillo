# Secure Samsung Vault

Sistema di sicurezza per Samsung Galaxy S: vault cifrato locale, lock/wipe/revoke remoto, compliance check.

## Requisiti

- Docker + Docker Compose
- Go 1.22+ (solo per sviluppo locale senza Docker)
- Android Studio Hedgehog+ (per l'app Android)
- JDK 17

## Avvio Backend

### Con Docker Compose (raccomandato)

```bash
# 1. Copia e configura il file .env
cp .env.example .env
# Modifica .env e imposta un ADMIN_TOKEN sicuro

# 2. Avvia tutto
docker compose up --build -d

# 3. Verifica che funzioni
curl http://localhost:8080/health
# Risposta: {"status":"ok"}
```

### Sviluppo locale (senza Docker per il backend)

```bash
# 1. Avvia solo PostgreSQL e Redis
docker compose up postgres redis -d

# 2. Esporta le variabili d'ambiente
export DATABASE_URL="postgres://vault:vault@localhost:5432/securevault?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export ADMIN_TOKEN="dev-admin-token"

# 3. Applica le migrations manualmente
psql "$DATABASE_URL" -f backend/migrations/001_create_devices.sql
psql "$DATABASE_URL" -f backend/migrations/002_create_device_commands.sql
psql "$DATABASE_URL" -f backend/migrations/003_create_audit_logs.sql

# 4. Avvia il server
cd backend
go run ./cmd/server/
```

## Test API

```bash
# Enrollment di un device
curl -X POST http://localhost:8080/api/v1/device/enroll \
  -H "Content-Type: application/json" \
  -d '{"device_model":"Galaxy S24","android_version":"14","app_version":"1.0.0"}'

# Heartbeat (usa il token ricevuto dall'enrollment)
curl -X POST http://localhost:8080/api/v1/device/heartbeat \
  -H "Authorization: Bearer <DEVICE_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"device_id":"<DEVICE_ID>","app_version":"1.0.0","battery_pct":85,"is_rooted":false}'

# Invio comando remoto (lock/wipe/revoke) - richiede admin token
curl -X POST http://localhost:8080/api/v1/admin/commands \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"device_id":"<DEVICE_ID>","type":"lock"}'

# Polling comandi pendenti
curl http://localhost:8080/api/v1/device/<DEVICE_ID>/commands \
  -H "Authorization: Bearer <DEVICE_TOKEN>"

# Ack comando
curl -X POST http://localhost:8080/api/v1/device/commands/<COMMAND_ID>/ack \
  -H "Authorization: Bearer <DEVICE_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"device_id":"<DEVICE_ID>","success":true,"message":"executed"}'
```

## App Android

### Setup

1. Apri la cartella `android/` in Android Studio
2. Sincronizza Gradle
3. Configura l'URL del backend in `app/build.gradle.kts`:
   - Emulatore: `http://10.0.2.2:8080/api/v1/` (default)
   - Device fisico sulla stessa rete: `http://<IP_PC>:8080/api/v1/`
4. Build & Run su device Samsung Galaxy S o emulatore

### Flusso utente

1. **Enrollment**: al primo avvio, il dispositivo si registra sul backend
2. **Setup PIN**: l'utente imposta un PIN a 6+ cifre
3. **Autenticazione**: ad ogni accesso, PIN o biometria
4. **Vault**: gestione dati cifrati (AES-256-GCM, chiavi in Android Keystore)
5. **Background**: heartbeat e polling comandi ogni 15 minuti via WorkManager
6. **Auto-lock**: vault si blocca dopo 5 minuti di inattività

### Comandi remoti

| Comando | Effetto |
|---------|---------|
| `lock`  | Blocca il vault, richiede ri-autenticazione |
| `wipe`  | Cancella tutti i dati cifrati e distrugge le chiavi |
| `revoke`| Wipe completo + de-registrazione del device |

## Struttura Progetto

```
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Configurazione
│   │   ├── database/        # PostgreSQL + Redis
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/       # Auth middleware
│   │   ├── models/          # Strutture dati
│   │   ├── repository/      # Data access layer
│   │   └── service/         # Business logic
│   ├── migrations/          # SQL migrations
│   └── Dockerfile
├── android/
│   └── app/src/main/java/com/gestionigrillo/securevault/
│       ├── compliance/       # Device compliance checks
│       ├── data/             # Database, network, repository impl
│       ├── domain/           # Models e interfacce
│       ├── presentation/     # Activities + ViewModels
│       ├── remotecommands/   # WorkManager workers
│       ├── securecommunications/  # Interfacce future
│       ├── vaultcrypto/      # Keystore + AES-GCM + PBKDF2
│       ├── vaultdb/          # Auto-lock
│       └── wipe/             # Vault destruction
└── docker-compose.yml
```
