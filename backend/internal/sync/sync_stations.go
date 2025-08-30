package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FlexibleNameField struct {
	Name string
}

func (f *FlexibleNameField) UnmarshalJSON(data []byte) error {
	// Handle "string"
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		f.Name = str
		return nil
	}

	// Handle object { "name": "..." }
	var obj struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	f.Name = obj.Name
	return nil
}

type SpaceDevsStation struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Status      FlexibleNameField `json:"status"`
	Type        FlexibleNameField `json:"type"`
	Orbit       FlexibleNameField `json:"orbit"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Founded     string            `json:"founded"`
}

type SpaceDevsResponseStations struct {
	Results []SpaceDevsStation `json:"results"`
}

func SyncStations() error {
	apiURL := "https://ll.thespacedevs.com/2.3.0/space_stations/?limit=30&format=json"

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch stations: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data SpaceDevsResponseStations
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("✅ Successfully fetched %d stations\n", len(data.Results))

	var errors []error
	for _, station := range data.Results {
		if err := createStationInPocketbase(station); err != nil {
			log.Printf("❌ Error inserting %s: %v", station.Name, err)
			errors = append(errors, err)
		} else {
			log.Printf("✅ Inserted station: %s", station.Name)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during sync (first error: %v)", len(errors), errors[0])
	}

	return nil
}

func createStationInPocketbase(station SpaceDevsStation) error {
	pbURL := "http://127.0.0.1:8080/api/collections/stations/records"

	payload := map[string]any{
		"api_id":      station.ID,
		"name":        station.Name,
		"status":      station.Status.Name,
		"type":        station.Type.Name,
		"orbit":       station.Orbit.Name,
		"url":         station.URL,
		"description": station.Description,
		"founded":     station.Founded,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", pbURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
	}

	return nil
}
