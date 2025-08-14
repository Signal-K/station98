package main

import (
	"log"
	"time"

	"github.com/signal-k/notifs/internal/config"
	"github.com/signal-k/notifs/internal/pbclient"
	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	log.Println("🚀 Backend is starting...")

	cfg := config.Load()

	client := pbclient.NewClient(cfg.PocketbaseURL)
	if err := client.Login(cfg.PocketbaseAdmin, cfg.PocketbasePassword); err != nil {
		log.Fatalf("Admin login failed: %v", err)
	}

	// Start the sync job in the background
	go func() {
		for {
			sync.SyncRocketLaunches(client)
			time.Sleep(6 * time.Hour)
		}
	}()

	// Prevent the app from exiting
	select {} // block forever
}
