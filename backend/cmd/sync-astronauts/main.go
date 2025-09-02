package main

import (
	"log"
	"os"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncAstronauts(); err != nil {
		log.Printf("❌ Failed to sync astronauts: %v", err)
		os.Exit(1)
	}
	
	log.Println("🎉 Astronaut sync completed successfully!")
}
