package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.MatchExpeditionLaunches(); err != nil {
		log.Fatalf("❌ Match failed: %v", err)
	}
}
