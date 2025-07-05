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

type SpaceDevsPayload struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Serial        string `json:"serial_number"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Nationalities []struct {
		Name string `json:"name"`
	} `json:"nationalities"`
	Orbit struct {
		Name string `json:"name"`
	} `json:"orbit"`
	Mass       float64 `json:"mass"`
	MassUnit   string  `json:"mass_unit"`
	Reusable   bool    `json:"reusable"`
	Updated    string  `json:"updated"`
	Spacecraft struct {
		Name string `json:"name"`
	} `json:"spacecraft"`
	SpacecraftConfig struct {
		Name string `json:"name"`
	} `json:"spacecraft_config"`
	SpacecraftStage struct {
		SpacecraftFlight struct {
			Flight string `json:"flight"`
		} `json:"spacecraft_flight"`
	} `json:"spacecraft_stage"`
	SpacecraftFlight string `json:"spacecraft_flight"`
	SpacecraftName   string `json:"spacecraft_name"`
}

type SpaceDevsPayloadResponse struct {
	Results []SpaceDevsPayload `json:"results"`
	Next    string             `json:"next"`
}

// SyncPayloads fetches all payloads and stores them in PocketBase
func SyncPayloads() error {
	fmt.Println("📦 Syncing payloads...")

	pageURL := "https://ll.thespacedevs.com/2.3.0/payloads/?limit=100"
	client := &http.Client{Timeout: 15 * time.Second}

	count := 0

	for pageURL != "" {
		resp, err := client.Get(pageURL)
		if err != nil {
			return fmt.Errorf("failed to fetch payloads: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		var data SpaceDevsPayloadResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		for _, payload := range data.Results {
			if err := createPayloadInPocketbase(payload); err != nil {
				log.Printf("❌ Failed to insert payload '%s': %v", payload.Name, err)
			} else {
				log.Printf("✅ Inserted payload: %s", payload.Name)
				count++
			}
		}

		pageURL = data.Next
	}

	fmt.Printf("✅ Completed syncing %d payloads.\n", count)
	return nil
}

func createPayloadInPocketbase(payload SpaceDevsPayload) error {
	pbURL := "http://127.0.0.1:8080/api/collections/payloads/records"

	nationality := ""
	if len(payload.Nationalities) > 0 {
		nationality = payload.Nationalities[0].Name
	}

	payloadBody := map[string]any{
		"api_id":            payload.ID,
		"name":              payload.Name,
		"slug":              payload.Slug,
		"description":       payload.Description,
		"serial_number":     payload.Serial,
		"nationality":       nationality,
		"orbit":             payload.Orbit.Name,
		"mass":              payload.Mass,
		"mass_unit":         payload.MassUnit,
		"reusable":          payload.Reusable,
		"spacecraft":        payload.Spacecraft.Name,
		"spacecraft_config": payload.SpacecraftConfig.Name,
		"spacecraft_flight": payload.SpacecraftStage.SpacecraftFlight.Flight,
		"updated_at":        payload.Updated,
	}

	body, err := json.Marshal(payloadBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", pbURL, bytes.NewBuffer(body))
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
		return fmt.Errorf("PocketBase HTTP %d: %s", res.StatusCode, string(body))
	}

	return nil
}
