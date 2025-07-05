package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncPayloads(); err != nil {
		log.Fatalf("sync failed: %v", err)
	}
}
