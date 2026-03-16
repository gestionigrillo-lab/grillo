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
	"github.com/gestionigrillo/secure-samsung-vault/internal/config"
	"github.com/gestionigrillo/secure-samsung-vault/internal/database"
	"github.com/gestionigrillo/secure-samsung-vault/internal/handlers"
	"github.com/gestionigrillo/secure-samsung-vault/internal/middleware"
	"github.com/gestionigrillo/secure-samsung-vault/internal/repository"
	"github.com/gestionigrillo/secure-samsung-vault/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Database connections
	pgPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pgPool.Close()

	redisClient, err := database.NewRedisClient(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	// Repositories
	deviceRepo := repository.NewDeviceRepository(pgPool)
	commandRepo := repository.NewCommandRepository(pgPool)
	auditRepo := repository.NewAuditRepository(pgPool)

	// Services
	deviceSvc := service.NewDeviceService(deviceRepo, auditRepo, redisClient)
	commandSvc := service.NewCommandService(commandRepo, deviceRepo, auditRepo)

	// Handlers
	deviceHandler := handlers.NewDeviceHandler(deviceSvc, commandSvc)
	adminHandler := handlers.NewAdminHandler(commandSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Device endpoints (enrollment is public, rest requires device token)
	r.Route("/api/v1", func(r chi.Router) {
		// Public: enrollment
		r.Post("/device/enroll", deviceHandler.Enroll)

		// Device-authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.DeviceAuth(deviceSvc.ValidateToken))
			r.Post("/device/heartbeat", deviceHandler.Heartbeat)
			r.Get("/device/{deviceId}/commands", deviceHandler.GetCommands)
			r.Post("/device/commands/{commandId}/ack", deviceHandler.AckCommand)
		})

		// Admin endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminAuth(cfg.AdminToken))
			r.Post("/admin/commands", adminHandler.CreateCommand)
		})
	})

	// Server
	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on %s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
