package main

import (
	"log"
	"os"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncMostRecentDockingLocation(); err != nil {
		log.Printf("❌ Docking location sync failed: %v", err)
		os.Exit(1)
	}
	log.Println("✅ Docking location sync completed successfully.")
}
