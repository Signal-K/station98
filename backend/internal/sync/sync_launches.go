package sync

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/signal-k/notifs/internal/pbclient"
)

type Launch struct {
	Name      string `json:"name"`
	Net       string `json:"net"`
	URL       string `json:"url"`
	LaunchPad struct {
		Name string `json:"name"`
	} `json:"pad"`
}

type LaunchAPIResponse struct {
	Results []Launch `json:"results"`
}

func SyncRocketLaunches(client *pbclient.Client) {
	log.Println("📡 Syncing launches from Launch Library...")

	resp, err := http.Get("https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=10")
	if err != nil {
		log.Printf("❌ Failed to fetch launches: %v", err)
		return
	}
	defer resp.Body.Close()

	var result LaunchAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("❌ Invalid JSON: %v", err)
		return
	}

	for _, l := range result.Results {
		launchTime, _ := time.Parse(time.RFC3339, l.Net)

		event := pbclient.Event{
			Title:       l.Name,
			Type:        "rocket_launch",
			Datetime:    launchTime.Format(time.RFC3339),
			Location:    l.LaunchPad.Name,
			SourceURL:   l.URL,
			Description: "Synced from Launch Library",
		}

		if err := client.CreateEvent(event); err != nil {
			log.Printf("❌ Could not insert: %s → %v", l.Name, err)
		} else {
			log.Printf("✅ Synced: %s", l.Name)
		}
	}
}
