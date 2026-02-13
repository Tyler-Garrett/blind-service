package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tyler-garrett/blind-service/internal/api"
	"github.com/tyler-garrett/blind-service/internal/api/handlers"
	"github.com/tyler-garrett/blind-service/internal/config"
	"github.com/tyler-garrett/blind-service/internal/repository"
	"github.com/tyler-garrett/blind-service/internal/service"
	"github.com/tyler-garrett/blind-service/pkg/storage"
)

const serviceName = "blind-service"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting %s...", serviceName)

	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Port: %s", cfg.Server.Port)

	// Initialize Azure Table Storage
	storageClient, err := storage.NewTableClient(cfg.Database.TableStorageConnectionString)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	ctx := context.Background()

	// Ensure table exists
	if err := storageClient.EnsureTable(ctx, cfg.Database.BlindsTableName); err != nil {
		log.Printf("Warning: failed to ensure table exists: %v", err)
	}

	// Initialize repository
	blindRepo := repository.NewBlindRepository(
		storageClient.GetTableClient(cfg.Database.BlindsTableName),
	)

	// Initialize service
	blindService := service.NewBlindService(blindRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	blindHandler := handlers.NewBlindHandler(blindService)

	// Setup router
	router := api.SetupRouter(cfg, healthHandler, blindHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
