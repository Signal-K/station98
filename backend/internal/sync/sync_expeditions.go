package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SpaceDevsExpedition struct {
	Name      string `json:"name"`
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
	URL       string `json:"url"`
	Station   struct {
		Name string `json:"name"`
	} `json:"spacestation"`
	Patch struct {
		Name  string `json:"name"`
		Image string `json:"image_url"`
	} `json:"mission_patch"`
	Crew []struct {
		Name string `json:"name"`
	} `json:"crew"`
}

type SpaceDevsResponseExpeditions struct {
	Results []SpaceDevsExpedition `json:"results"`
}

func SyncExpeditions() error {
	resp, err := http.Get("https://ll.thespacedevs.com/2.3.0/expeditions/?limit=30&ordering=-start")
	if err != nil {
		return fmt.Errorf("failed to fetch expeditions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read expedition response: %w", err)
	}

	var data SpaceDevsResponseExpeditions
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse expedition JSON: %w", err)
	}

	for _, exp := range data.Results {
		if err := createExpeditionInPocketbase(exp); err != nil {
			fmt.Printf("❌ Failed to insert %s: %v\n", exp.Name, err)
		} else {
			fmt.Printf("✅ Inserted expedition: %s\n", exp.Name)
		}
	}

	return nil
}

func createExpeditionInPocketbase(exp SpaceDevsExpedition) error {
	stationId, err := findStationIdByName(exp.Station.Name)
	if err != nil {
		return fmt.Errorf("station lookup failed: %w", err)
	}

	crewIds := []string{}
	for _, member := range exp.Crew {
		id, err := findAstronautIdByName(member.Name)
		if err == nil && id != "" {
			crewIds = append(crewIds, id)
		}
	}

	payload := map[string]any{
		"name":       exp.Name,
		"start_date": exp.StartDate,
		"end_date":   exp.EndDate,
		"station":    stationId,
		"url":        exp.URL,
		"patches":    exp.Patch.Image,
		"crew":       crewIds,
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/collections/expeditions/records", bytes.NewBuffer(jsonData))
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

func findStationIdByName(name string) (string, error) {
	encodedName := url.QueryEscape(fmt.Sprintf(`name="%s"`, name))
	query := fmt.Sprintf("http://127.0.0.1:8080/api/collections/stations/records?filter=%s", encodedName)

	resp, err := http.Get(query)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if len(result.Items) == 0 {
		return "", fmt.Errorf("station %s not found", name)
	}

	return result.Items[0].ID, nil
}

func findAstronautIdByName(name string) (string, error) {
	query := fmt.Sprintf("http://127.0.0.1:8080/api/collections/astronauts/records?filter=name='%s'", name)
	resp, err := http.Get(query)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if len(result.Items) == 0 {
		return "", fmt.Errorf("astronaut %s not found", name)
	}

	return result.Items[0].ID, nil
}
