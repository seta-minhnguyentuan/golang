package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asset-service/internal/config"
	"asset-service/internal/database"
	httpserver "asset-service/internal/http"
	"asset-service/internal/repository"
	"asset-service/internal/services"
)

func main() {
	dbCfg := config.LoadDB()
	srvCfg := config.LoadServerConfig()

	db, err := database.Connect(*dbCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to database")

	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migrations completed")

	folderRepo := repository.NewFolderRepository(db)
	folderSvc := services.NewFolderService(folderRepo)

	engine := httpserver.NewRouter(httpserver.RouterDeps{
		FolderService: folderSvc,
	})

	srv := &http.Server{
		Addr:           ":" + srvCfg.Port,
		Handler:        engine,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("HTTP server listening on :%s", srvCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	_ = srv.Close()
}
