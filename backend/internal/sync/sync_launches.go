package sync

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/signal-k/notifs/internal/pbclient"
)

type Agency struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Abbrev   string `json:"abbrev"`
	Type     string `json:"type"`
	Country  string `json:"country_code"`
	Desc     string `json:"description"`
	Founded  string `json:"founding_year"`
	LogoURL  string `json:"logo_url"`
	ImageURL string `json:"image_url"`
	WikiURL  string `json:"wiki_url"`
	InfoURL  string `json:"info_url"`
}

type Launch struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Net                   string `json:"net"`
	URL                   string `json:"url"`
	LaunchServiceProvider Agency `json:"launch_service_provider"`
	Pad                   struct {
		Name string `json:"name"`
	} `json:"pad"`
}

type LaunchAPIResponse struct {
	Results []Launch `json:"results"`
}

func getThrottleDelay(body string) time.Duration {
	// Try to extract seconds from the error message
	if idx := strings.Index(body, "Expected available in "); idx != -1 {
		rest := body[idx+21:]
		secStr := ""
		for _, c := range rest {
			if c >= '0' && c <= '9' {
				secStr += string(c)
			} else {
				break
			}
		}
		if secStr != "" {
			secs, err := strconv.Atoi(secStr)
			if err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}
	return 60 * time.Second // fallback
}

func SyncLaunchProvidersAndEvents(client *pbclient.Client) {
	log.Println("📡 Fetching upcoming launches...")

	resp, err := http.Get("https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=10")
	if err != nil {
		log.Printf("❌ Failed to fetch launches: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		body, _ := io.ReadAll(resp.Body)
		delay := getThrottleDelay(string(body))
		log.Printf("❌ Launches: HTTP 429: %s. Retrying in %v...", string(body), delay)
		time.Sleep(delay)
		SyncLaunchProvidersAndEvents(client)
		return
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Launches: HTTP %d: %s", resp.StatusCode, string(body))
		return
	}

	var result LaunchAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("❌ Invalid JSON: %v", err)
		return
	}

	providerIDs := make(map[int]Agency)
	for _, l := range result.Results {
		provider := l.LaunchServiceProvider
		if provider.ID != 0 && provider.Name != "" {
			providerIDs[provider.ID] = provider
		}
	}

	// Ensure each provider exists before creating events
	for _, provider := range providerIDs {
		existing, err := client.FindRecordByField("launch_providers", "spacedevs_id", provider.ID)
		if err == nil && existing != nil {
			log.Printf("⏭️ Provider already exists: %s", provider.Name)
			continue
		}
		_, err = client.CreateRecord("launch_providers", map[string]interface{}{
			"spacedevs_id":  provider.ID,
			"name":          provider.Name,
			"abbrev":        provider.Abbrev,
			"type":          provider.Type,
			"country_code":  provider.Country,
			"description":   provider.Desc,
			"founding_year": provider.Founded,
			"logo_url":      provider.LogoURL,
			"image_url":     provider.ImageURL,
			"wiki_url":      provider.WikiURL,
			"info_url":      provider.InfoURL,
		})
		if err != nil {
			log.Printf("❌ Failed to insert provider %s: %v", provider.Name, err)
		} else {
			log.Printf("✅ Synced provider: %s", provider.Name)
		}
	}

	log.Println("✅ Provider sync complete.")
	log.Println("📡 Syncing launches/events...")

	for _, l := range result.Results {
		// Check for existing event by title (unique field)
		existingByTitle, errTitle := client.FindRecordByField("events", "title", l.Name)
		if errTitle == nil && existingByTitle != nil {
			log.Printf("⚠️ Launch %s already exists (by title), skipping.", l.Name)
			continue
		}

		launchTime, _ := time.Parse(time.RFC3339, l.Net)

		// Find the PocketBase provider record for this launch
		providerRecord, err := client.FindRecordByField("launch_providers", "spacedevs_id", l.LaunchServiceProvider.ID)
		if err != nil || providerRecord == nil {
			log.Printf("❌ Could not find provider for launch %s (id %d)", l.Name, l.LaunchServiceProvider.ID)
			continue
		}
		providerPBID := (*providerRecord)["id"].(string)

		event := map[string]interface{}{
			"title":        l.Name,
			"type":         "rocket_launch",
			"datetime":     launchTime.Format(time.RFC3339),
			"location":     l.Pad.Name,
			"source_url":   l.URL,
			"description":  "Synced from Launch Library",
			"spacedevs_id": providerPBID, // relation to launch_providers
		}

		_, err = client.CreateRecord("events", event)
		if err != nil {
			log.Printf("❌ Could not insert: %s → %v", l.Name, err)
		} else {
			log.Printf("✅ Synced: %s", l.Name)
		}
	}
}
