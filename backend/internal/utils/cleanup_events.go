package utils

import (
	"log"
	"strings"

	"github.com/signal-k/notifs/internal/pbclient"
)

// RemoveDuplicateEvents finds events with duplicate titles and removes all but one.
func RemoveDuplicateEvents(client *pbclient.Client) error {
	// Fetch all records (handle pagination)
	var allEvents []map[string]interface{}
	page := 1
	for {
		events, err := client.ListRecordsWithPage("events", page, 100)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			break
		}
		allEvents = append(allEvents, events...)
		page++
	}

	seen := make(map[string]string)
	duplicates := []string{}

	for _, event := range allEvents {
		rawTitle, ok := event["title"]
		if !ok {
			continue
		}
		title, ok := rawTitle.(string)
		if !ok {
			continue
		}
		normalized := strings.ToLower(strings.TrimSpace(title))
		id, ok := event["id"].(string)
		if !ok {
			continue
		}
		if prevID, exists := seen[normalized]; exists {
			duplicates = append(duplicates, id)
			log.Printf("Duplicate event found: %s (keeping %s, removing %s)", title, prevID, id)
		} else {
			seen[normalized] = id
		}
	}

	for _, id := range duplicates {
		if err := client.DeleteRecord("events", id); err != nil {
			log.Printf("Failed to delete duplicate event %s: %v", id, err)
		} else {
			log.Printf("Deleted duplicate event %s", id)
		}
	}

	log.Printf("Cleanup complete. Removed %d duplicates.", len(duplicates))
	return nil
}
