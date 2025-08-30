package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	log.Println("🛰️ Fetching docking locations...")

	if err := sync.SyncDockingLocations(); err != nil {
		log.Fatalf("Sync failed: %v", err)
	}

	log.Println("✅ Docking location sync completed.")
}
