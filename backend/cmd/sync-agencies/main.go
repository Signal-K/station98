package main

import (
	"log"

	"github.com/signal-k/notifs/internal/sync"
)

func main() {
	if err := sync.SyncAgencies(); err != nil {
		log.Fatal(err)
	}
}
