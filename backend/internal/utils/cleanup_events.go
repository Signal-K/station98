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

	log.Printf("Cleanup complete. Removed %d duplicates from events.", len(duplicates))

	// --- Remove duplicates in missions collection ---
	var allMissions []map[string]interface{}
	page = 1
	for {
		missions, err := client.ListRecordsWithPage("missions", page, 100)
		if err != nil {
			return err
		}
		if len(missions) == 0 {
			break
		}
		allMissions = append(allMissions, missions...)
		page++
	}

	seenM := make(map[string]string)
	duplicatesM := []string{}

	for _, mission := range allMissions {
		rawName, ok := mission["name"]
		if !ok {
			continue
		}
		name, ok := rawName.(string)
		if !ok {
			continue
		}
		normalized := strings.ToLower(strings.TrimSpace(name))
		id, ok := mission["id"].(string)
		if !ok {
			continue
		}
		if prevID, exists := seenM[normalized]; exists {
			duplicatesM = append(duplicatesM, id)
			log.Printf("Duplicate mission found: %s (keeping %s, removing %s)", name, prevID, id)
		} else {
			seenM[normalized] = id
		}
	}

	for _, id := range duplicatesM {
		if err := client.DeleteRecord("missions", id); err != nil {
			log.Printf("Failed to delete duplicate mission %s: %v", id, err)
		} else {
			log.Printf("Deleted duplicate mission %s", id)
		}
	}

	log.Printf("Cleanup complete. Removed %d duplicates from missions.", len(duplicatesM))

	return nil
}
