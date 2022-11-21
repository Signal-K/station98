package main

import (
	"log"
	"os"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncPrograms(); err != nil {
		log.Printf("❌ Failed to sync programs: %v", err)
		os.Exit(1)
	}

	log.Println("🎉 Program sync completed successfully!")
}
