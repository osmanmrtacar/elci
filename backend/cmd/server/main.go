package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/osmanmertacar/sosyal/backend/internal/api"
	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting TikTok Content Publisher API...")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Server will listen on: %s", cfg.GetServerAddress())

	// Initialize database
	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Setup HTTP router
	router := api.SetupRouter(cfg, db)

	// Setup graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server started successfully on %s", addr)
	log.Printf("Health check available at: http://%s/health", cfg.GetServerAddress())
	log.Printf("API base URL: http://%s/api/v1", cfg.GetServerAddress())

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
