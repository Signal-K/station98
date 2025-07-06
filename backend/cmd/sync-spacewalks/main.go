package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncSpacewalks(); err != nil {
		log.Fatalf("❌ Spacewalk sync failed: %v", err)
	}
}
