package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/handlers"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/loader"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/repository"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/service"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	repo := repository.NewMemoryRepository()
	svc := service.NewItemService(repo)

	// Preload items from CSV
	csvLoader := loader.NewCSVLoader(repo)
	csvPath := os.Getenv("CSV_DATA_PATH")
	if csvPath == "" {
		csvPath = "data/items.csv"
	}

	if _, err := os.Stat(csvPath); err == nil {
		ctx := context.Background()
		count, err := csvLoader.LoadFromFile(ctx, csvPath)
		if err != nil {
			log.Printf("Warning: Failed to preload items from CSV: %v", err)
		} else {
			log.Printf("Successfully preloaded %d items from %s", count, csvPath)
		}
	} else {
		log.Printf("CSV file not found at %s, starting with empty repository", csvPath)
	}

	router := handlers.SetupRoutes(svc)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
