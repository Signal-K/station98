package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SpaceDevsDockingLocationResponse struct {
	Results []SpaceDevsDockingLocation `json:"results"`
}

type SpaceDevsDockingLocation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Spacestation struct {
		ID    int    `json:"id"`
		URL   string `json:"url"`
		Name  string `json:"name"`
		Image struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			ImageURL     string `json:"image_url"`
			ThumbnailURL string `json:"thumbnail_url"`
			Credit       string `json:"credit"`
			License      struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Priority int    `json:"priority"`
				Link     string `json:"link"`
			} `json:"license"`
			SingleUse bool `json:"single_use"`
		} `json:"image"`
	} `json:"spacestation"`
	Spacecraft *struct {
		Name string `json:"name"`
	} `json:"spacecraft"`
	Payload *struct {
		Name string `json:"name"`
	} `json:"payload"`
}

// SyncMostRecentDockingLocation fetches and stores the most recent docking location
func SyncMostRecentDockingLocation() error {
	fmt.Println("🛰️ Fetching most recent docking location...")

	resp, err := http.Get("https://ll.thespacedevs.com/2.3.0/config/docking_locations/?limit=1&ordering=-id&format=json")
	if err != nil {
		return fmt.Errorf("failed to fetch docking location: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data SpaceDevsDockingLocationResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse docking location response: %w", err)
	}
	if len(data.Results) == 0 {
		return fmt.Errorf("no docking location results found")
	}

	docking := data.Results[0]
	if err := createDockingLocationInPocketbase(docking); err != nil {
		return fmt.Errorf("failed to insert docking location: %w", err)
	}

	fmt.Println("✅ Most Recent Docking Location:")
	fmt.Printf("- Name: %s\n", docking.Name)
	fmt.Printf("- Station: %s\n", docking.Spacestation.Name)
	fmt.Printf("  • Image: %s\n", docking.Spacestation.Image.ImageURL)
	fmt.Printf("  • Credit: %s\n", docking.Spacestation.Image.Credit)
	fmt.Printf("  • License: %s (%s)\n", docking.Spacestation.Image.License.Name, docking.Spacestation.Image.License.Link)
	if docking.Spacecraft != nil {
		fmt.Printf("- Spacecraft: %s\n", docking.Spacecraft.Name)
	} else {
		fmt.Println("- Spacecraft: None")
	}
	if docking.Payload != nil {
		fmt.Printf("- Payload: %s\n", docking.Payload.Name)
	} else {
		fmt.Println("- Payload: None")
	}

	return nil
}

func createDockingLocationInPocketbase(d SpaceDevsDockingLocation) error {
	// Replace with your actual Pocketbase station record mapping logic
	stationPbID, err := findStationPocketbaseID(d.Spacestation.ID)
	if err != nil {
		return fmt.Errorf("could not resolve station: %w", err)
	}

	payload := map[string]interface{}{
		"api_id":               d.ID,
		"name":                 d.Name,
		"station":              stationPbID,
		"spacecraft_name":      derefString(d.Spacecraft),
		"payload_name":         derefString(d.Payload),
		"station_image_url":    d.Spacestation.Image.ImageURL,
		"station_image_credit": d.Spacestation.Image.Credit,
		"station_license_name": d.Spacestation.Image.License.Name,
		"station_license_url":  d.Spacestation.Image.License.Link,
		"api_url":              fmt.Sprintf("https://ll.thespacedevs.com/2.3.0/config/docking_locations/%d/", d.ID),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/collections/docking_locations/records", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
	}

	return nil
}

// Helper to safely dereference name field from *structs
func derefString(ptr interface{}) string {
	switch v := ptr.(type) {
	case *struct{ Name string }:
		if v != nil {
			return v.Name
		}
	case nil:
		return ""
	}
	return ""
}

// Replace with real lookup (from Pocketbase stations by api_id)
func findStationPocketbaseID(spaceDevsStationID int) (string, error) {
	// TEMPORARY placeholder — replace with Pocketbase fetch (GET /collections/stations)
	// that maps api_id -> Pocketbase ID
	switch spaceDevsStationID {
	case 4:
		return "RECORD_ID_FOR_ISS", nil
	default:
		return "", fmt.Errorf("no known mapping for station ID %d", spaceDevsStationID)
	}
}
