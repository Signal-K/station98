package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	log.Println("🛰️ Fetching space stations...")

	if err := sync.SyncStations(); err != nil {
		log.Fatalf("Sync failed: %v", err)
	}

	log.Println("✅ Sync completed successfully.")
}
