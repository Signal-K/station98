package main

import (
	"log"

	"github.com/signal-k/notifs/internal/utils"
)

func runVidURLs() {
	if err := utils.ListLaunchesWithVidURLs(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
