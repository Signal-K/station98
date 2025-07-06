package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncDockingEvents(); err != nil {
		log.Fatalf("🚨 Sync failed: %v", err)
	}
}
