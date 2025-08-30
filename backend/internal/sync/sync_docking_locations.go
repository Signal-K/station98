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

type SpaceDevsDockingLocation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Spacestation struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Image struct {
			ImageURL string `json:"image_url"`
			Credit   string `json:"credit"`
			License  struct {
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"license"`
		} `json:"image"`
	} `json:"spacestation"`
}

type SpaceDevsDockingLocationResponse struct {
	Results []SpaceDevsDockingLocation `json:"results"`
}

func SyncDockingLocations() error {
	apiURL := "https://ll.thespacedevs.com/2.3.0/config/docking_locations/?limit=30&format=json"
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch docking locations: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data SpaceDevsDockingLocationResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("✅ Fetched %d docking locations\n", len(data.Results))

	for _, loc := range data.Results {
		if err := insertDockingLocation(loc); err != nil {
			log.Printf("❌ Failed to insert docking location %s: %v", loc.Name, err)
		} else {
			log.Printf("✅ Inserted: %s (Station: %s)", loc.Name, loc.Spacestation.Name)
		}
	}

	return nil
}

func insertDockingLocation(loc SpaceDevsDockingLocation) error {
	// Find the corresponding station in Pocketbase by name
	stationID, err := findStationIDByName(loc.Spacestation.Name)
	if err != nil {
		return fmt.Errorf("station lookup failed: %w", err)
	}

	pbURL := "http://127.0.0.1:8080/api/collections/docking_locations/records"
	payload := map[string]any{
		"api_id":       loc.ID,
		"name":         loc.Name,
		"station":      stationID,
		"station_url":  loc.Spacestation.URL,
		"image_url":    loc.Spacestation.Image.ImageURL,
		"image_credit": loc.Spacestation.Image.Credit,
		"license_name": loc.Spacestation.Image.License.Name,
		"license_url":  loc.Spacestation.Image.License.Link,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", pbURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("POST to Pocketbase: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Pocketbase error %d: %s", res.StatusCode, body)
	}

	return nil
}

func findStationIDByName(stationName string) (string, error) {
	query := url.QueryEscape(fmt.Sprintf("name='%s'", stationName))
	apiURL := fmt.Sprintf("http://127.0.0.1:8080/api/collections/stations/records?filter=%s", query)

	res, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("station lookup HTTP error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("station lookup failed: HTTP %d: %s", res.StatusCode, body)
	}

	var response struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode station lookup: %w", err)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("station not found for name: %s", stationName)
	}

	return response.Items[0].ID, nil
}
