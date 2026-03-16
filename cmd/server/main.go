package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/secure-samsung-vault/grillo/internal/config"
	"github.com/secure-samsung-vault/grillo/internal/database"
	"github.com/secure-samsung-vault/grillo/internal/handlers"
	"github.com/secure-samsung-vault/grillo/internal/middleware"
	"github.com/secure-samsung-vault/grillo/internal/repository"
	"github.com/secure-samsung-vault/grillo/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	redisClient, err := database.NewRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	// Repositories
	deviceRepo := repository.NewDeviceRepository(pool)
	commandRepo := repository.NewCommandRepository(pool)
	auditRepo := repository.NewAuditRepository(pool)

	// Services
	auditSvc := service.NewAuditService(auditRepo)
	deviceSvc := service.NewDeviceService(deviceRepo, auditSvc)
	commandSvc := service.NewCommandService(commandRepo, deviceRepo, auditSvc)

	// Handlers
	deviceHandler := handlers.NewDeviceHandler(deviceSvc, commandSvc)
	commandHandler := handlers.NewCommandHandler(commandSvc)
	adminHandler := handlers.NewAdminHandler(commandSvc)

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestLogging)
	r.Use(chimw.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Device endpoints (protected by device auth)
		r.Route("/device", func(r chi.Router) {
			r.Use(middleware.DeviceAuth())
			r.Post("/enroll", deviceHandler.Enroll)
			r.Post("/heartbeat", deviceHandler.Heartbeat)
			r.Get("/{deviceId}/commands", deviceHandler.GetPendingCommands)
			r.Post("/commands/{commandId}/ack", commandHandler.AcknowledgeCommand)
		})

		// Admin endpoints (protected by admin auth)
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminAuth(cfg.AdminToken))
			r.Post("/commands", adminHandler.IssueCommand)
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-done
	log.Println("server shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}

	log.Println("server stopped")
}
