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
	MissionPatches []struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
	} `json:"mission_patches"`
	Crew []struct {
		Astronaut struct {
			Name string `json:"name"`
		} `json:"astronaut"`
	} `json:"crew"`
}

type SpaceDevsResponseExpeditions struct {
	Results []SpaceDevsExpedition `json:"results"`
}

func SyncExpeditions() error {
	fmt.Println("🚀 Fetching expeditions...")

	resp, err := http.Get("https://ll.thespacedevs.com/2.3.0/expeditions/?limit=30&ordering=-start&mode=detailed")
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
	var stationId string
	var err error

	if exp.Station.Name != "" {
		stationId, err = findStationIdByName(exp.Station.Name)
		if err != nil {
			fmt.Printf("⚠️  Station not found for expedition %s: %s\n", exp.Name, exp.Station.Name)
			stationId = ""
		}
	}

	crewIds := []string{}
	for _, member := range exp.Crew {
		name := member.Astronaut.Name
		id, err := findAstronautIdByName(name)
		if err == nil && id != "" {
			crewIds = append(crewIds, id)
		}
	}

	// Pick first mission patch if any
	patchImage := ""
	if len(exp.MissionPatches) > 0 {
		patchImage = exp.MissionPatches[0].ImageURL
	}
	fmt.Printf("🎖️  Expedition '%s' patch image: %s\n", exp.Name, patchImage)

	payload := map[string]any{
		"name":       exp.Name,
		"start_date": exp.StartDate,
		"url":        exp.URL,
		"crew":       crewIds,
	}

	if exp.EndDate != "" {
		payload["end_date"] = exp.EndDate
	}
	if stationId != "" {
		payload["station"] = stationId
	}
	if patchImage != "" {
		payload["patches"] = patchImage
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal expedition: %w", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/collections/expeditions/records", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to build POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	fmt.Printf("✅ Success for %s\n", exp.Name)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("station %s not found", name)
	}

	return result.Items[0].ID, nil
}

func findAstronautIdByName(name string) (string, error) {
	encodedName := url.QueryEscape(fmt.Sprintf(`name="%s"`, name))
	query := fmt.Sprintf("http://127.0.0.1:8080/api/collections/astronauts/records?filter=%s", encodedName)

	resp, err := http.Get(query)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("astronaut %s not found", name)
	}

	return result.Items[0].ID, nil
}
