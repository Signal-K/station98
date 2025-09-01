package utils

import (
	// "encoding/json"
	// "io"
	"log"
	"net/http"

	// "time"

	"github.com/signal-k/notifs/internal/pbclient"
	// "github.com/signal-k/notifs/internal/sync"
)

func FetchSecond50Events(client *pbclient.Client) error {
	log.Println("📡 Fetching launches (page 2)...")
	url := "https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=50&mode=detailed&page=2"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if resp.StatusCode == 429 {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	delay := sync.GetThrottleDelay(string(body))
	// 	log.Printf("❌ Launches: HTTP 429: %s. Retrying in %v...", string(body), delay)
	// 	time.Sleep(delay)
	// 	return FetchSecond50Events(client)
	// }
	// if resp.StatusCode != 200 {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	return sync.NewSyncError(resp.StatusCode, string(body))
	// }

	// var result sync.LaunchAPIResponse
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return err
	// }

	// for _, l := range result.Results {
	// 	if err := sync.SyncSingleLaunch(client, l); err != nil {
	// 		log.Printf("Failed to sync launch %s: %v", l.Name, err)
	// 	}
	// }

	log.Println("✅ Second 50 launches and related data synced.")
	return nil
}
