package main

import (
	"log"
	"os"

	"github.com/signal-k/notifs/internal/pbclient"
	"github.com/signal-k/notifs/internal/utils"
)

func main() {
	baseURL := "http://localhost:8080"
	adminEmail := os.Getenv("POCKETBASE_ADMIN_EMAIL")
	adminPassword := os.Getenv("POCKETBASE_ADMIN_PASSWORD")

	client := pbclient.NewClient(baseURL)
	if err := client.Login(adminEmail, adminPassword); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	if err := utils.RemoveDuplicateEvents(client); err != nil {
		log.Fatalf("Cleanup failed: %v", err)
	}
}
