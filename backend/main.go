package main

import (
	"log"
	"time"

	"github.com/signal-k/notifs/internal/config"
	"github.com/signal-k/notifs/internal/pbclient"
	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	log.Println("🚀 Go starting...")

	cfg := config.Load()

	client := pbclient.NewClient(cfg.PocketbaseURL)

	time.Sleep(2 * time.Second)

	// Retry admin login until PocketBase is ready
	var err error
	for i := 0; i < 30; i++ {
		err = client.Login(cfg.PocketbaseAdmin, cfg.PocketbasePassword)
		if err == nil {
			break
		}
		log.Printf("Waiting for PocketBase to be ready (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Admin login failed after retries: %v", err)
	}

	go func() {
		for {
			sync.SyncLaunchProvidersAndEvents(client)
			time.Sleep(6 * time.Hour)
		}
	}()

	select {}
}
