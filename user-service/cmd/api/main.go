package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"user-service/internal/app"
	"user-service/internal/config"
	"user-service/internal/database"
	httpserver "user-service/internal/http"
)

func main() {
	// Load configurations
	dbCfg := config.LoadDB()
	srvCfg := config.LoadServerConfig()
	kafkaCfg := config.LoadKafkaConfig()

	// Connect to database
	db, err := database.Connect(*dbCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to database")

	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations completed")

	// Wire up dependencies with Kafka support
	components := app.Wire(kafkaCfg)
	log.Println("✅ Components wired with Kafka support")

	engine := httpserver.NewRouter(httpserver.RouterDeps{
		UserService: components.Users,
		TeamService: components.Teams,
	})

	srv := &http.Server{
		Addr:           ":" + srvCfg.Port,
		Handler:        engine,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 HTTP server listening on :%s", srvCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()

	// Start Kafka consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 Starting Kafka consumer for topic: %s", kafkaCfg.KafkaTopicTeamActivity)
		if err := components.Consumer.Run(ctx); err != nil {
			log.Printf("❌ Kafka consumer error: %v", err)
		}
		log.Println("✅ Kafka consumer stopped")
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down...")

	// Cancel context to stop Kafka consumer
	cancel()

	// Shutdown HTTP server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v", err)
	}

	// Close Kafka components
	if err := components.Consumer.Close(); err != nil {
		log.Printf("❌ Error closing Kafka consumer: %v", err)
	}
	if err := components.Producer.Close(); err != nil {
		log.Printf("❌ Error closing Kafka producer: %v", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("✅ Shutdown complete")
}
