package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// SpaceDevsSpacewalk defines the structure of a spacewalk from SpaceDevs API
type SpaceDevsSpacewalk struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	Location  string `json:"location"`
	StartTime string `json:"start"` // ISO8601
	EndTime   string `json:"end"`
	Duration  string `json:"duration"`

	Expedition *struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"expedition"`

	Event *struct {
		ID int `json:"id"`
	} `json:"event"`
}

// SpaceDevsSpacewalkResponse wraps the API results
type SpaceDevsSpacewalkResponse struct {
	Count   int                  `json:"count"`
	Results []SpaceDevsSpacewalk `json:"results"`
}

// SyncSpacewalks fetches and syncs all spacewalks from TSD into Pocketbase
func SyncSpacewalks() error {
	fmt.Println("🚀 Syncing spacewalks from TSD...")

	url := "https://ll.thespacedevs.com/2.3.0/spacewalks/?limit=5000&format=json&ordering=-start"
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch spacewalks: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read spacewalk response: %w", err)
	}

	var data SpaceDevsSpacewalkResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to unmarshal spacewalk JSON: %w", err)
	}

	success := 0
	skipped := 0

	for _, sw := range data.Results {
		if err := createSpacewalkInPocketbase(sw); err != nil {
			if err.Error() == "skip" {
				skipped++
				continue
			}
			log.Printf("❌ Failed to insert spacewalk: %s — %v", sw.Name, err)
		} else {
			success++
			log.Printf("✅ Inserted spacewalk: %s", sw.Name)
		}
	}

	fmt.Printf("✅ Done syncing spacewalks: %d inserted, %d skipped\n", success, skipped)
	return nil
}

// createSpacewalkInPocketbase inserts a spacewalk into Pocketbase
func createSpacewalkInPocketbase(sw SpaceDevsSpacewalk) error {
	pbURL := "http://127.0.0.1:8080/api/collections/spacewalks/records"

	// Resolve expedition Pocketbase ID
	var expeditionPBID string
	if sw.Expedition != nil {
		id, err := getPocketbaseExpeditionID(sw.Expedition.ID)
		if err != nil {
			return fmt.Errorf("lookup expedition: %w", err)
		}
		if id == "" {
			log.Printf("⚠️ Skipping spacewalk %q: expedition %d not found in PB", sw.Name, sw.Expedition.ID)
			return fmt.Errorf("skip")
		}
		expeditionPBID = id
	}

	payload := map[string]interface{}{
		"api_id":     sw.ID,
		"name":       sw.Name,
		"slug":       sw.Slug,
		"url":        sw.URL,
		"location":   sw.Location,
		"start_time": sw.StartTime,
		"end_time":   sw.EndTime,
		"duration":   sw.Duration,
		"event_id":   sw.Event.ID,
	}

	// Only include expedition_name if Expedition is not nil
	if sw.Expedition != nil {
		payload["expedition_name"] = sw.Expedition.Name
	}

	if expeditionPBID != "" {
		payload["expedition"] = expeditionPBID
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", pbURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PB HTTP %d: %s", resp.StatusCode, body)
	}

	return nil
}

// getPocketbaseExpeditionID resolves a TSD API expedition ID to its PB record ID
func getPocketbaseExpeditionID(apiID int) (string, error) {
	filter := url.QueryEscape(fmt.Sprintf("api_id=%d", apiID))
	endpoint := fmt.Sprintf("http://127.0.0.1:8080/api/collections/expeditions/records?filter=%s", filter)

	res, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("lookup failed (status %d)", res.StatusCode)
	}

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", nil // Not found
	}
	return result.Items[0].ID, nil
}
