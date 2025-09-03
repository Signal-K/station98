package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/signal-k/notifs/internal/config"
	"github.com/signal-k/notifs/internal/pbclient"
	"github.com/signal-k/notifs/internal/utils"
)

func main() {
	// Define command line flags
	var (
		cleanupEvents = flag.Bool("cleanup-events", false, "Remove duplicate events from the database")
		help          = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Check if any action flag was provided
	if !*cleanupEvents {
		log.Println("No action specified. Use -help to see available options.")
		os.Exit(1)
	}

	// Load configuration using the same config package as main app
	cfg := config.Load()

	// Create client and login with retry logic (similar to main app)
	client := pbclient.NewClient(cfg.PocketbaseURL)

	log.Println("🔧 Utility starting...")
	time.Sleep(1 * time.Second)

	// Retry admin login until PocketBase is ready
	var err error
	for i := 0; i < 10; i++ {
		err = client.Login(cfg.PocketbaseAdmin, cfg.PocketbasePassword)
		if err == nil {
			break
		}
		log.Printf("Waiting for PocketBase to be ready (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Admin login failed after retries: %v", err)
	}

	// Execute the requested action
	if *cleanupEvents {
		log.Println("Starting duplicate events cleanup...")
		if err := utils.RemoveDuplicateEvents(client); err != nil {
			log.Fatalf("Cleanup failed: %v", err)
		}
		log.Println("✅ Cleanup completed successfully!")
	}
}

func printHelp() {
	log.Println("Space Notifications Utility Tool")
	log.Println("")
	log.Println("Usage: go run cmd/utils-main.go [flags]")
	log.Println("")
	log.Println("Available flags:")
	log.Println("  -cleanup-events    Remove duplicate events from the database")
	log.Println("  -help             Show this help message")
	log.Println("")
	log.Println("Environment variables required:")
	log.Println("  PB_URL            PocketBase URL (e.g., http://localhost:8080)")
	log.Println("  PB_ADMIN_EMAIL    PocketBase admin email")
	log.Println("  PB_ADMIN_PASSWORD PocketBase admin password")
	log.Println("")
	log.Println("Examples:")
	log.Println("  go run cmd/utils-main.go -cleanup-events")
	log.Println("  go run cmd/utils-main.go -help")
}
