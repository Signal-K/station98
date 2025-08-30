package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	log.Println("🧭 Fetching space expeditions...")
	if err := sync.SyncExpeditions(); err != nil {
log.Fatalf("Sync failed: failed to parse expedition JSON: %v", err)
	}
}
