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

type SpaceDevsProgram struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        struct {
		Name string `json:"name"`
	} `json:"type"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	InfoURL   string `json:"info_url"`
	WikiURL   string `json:"wiki_url"`
	URL       string `json:"url"`
	Image     struct {
		ImageURL     string `json:"image_url"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"image"`
}

type SpaceDevsResponseProgram struct {
	Results []SpaceDevsProgram `json:"results"`
}

func SyncPrograms() error {
	fmt.Println("🚀 Fetching latest space programs...")

	apiURL := "https://ll.thespacedevs.com/2.2.0/program/?limit=3000&ordering=-start_date&mode=normal"
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch programs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data SpaceDevsResponseProgram
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("✅ Successfully fetched %d programs\n", len(data.Results))

	var errors []error
	for _, prog := range data.Results {
		if err := createProgramInPocketbase(prog); err != nil {
			log.Printf("❌ Error inserting %s: %v", prog.Name, err)
			errors = append(errors, fmt.Errorf("failed to insert %s: %w", prog.Name, err))
		} else {
			log.Printf("✅ Inserted program: %s", prog.Name)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during sync: %v", len(errors), errors[0])
	}

	return nil
}

func createProgramInPocketbase(prog SpaceDevsProgram) error {
	pbURL := "http://127.0.0.1:8080/api/collections/programs/records"

	payload := map[string]any{
		"api_id":          prog.ID,
		"name":            prog.Name,
		"type":            prog.Type.Name,
		"description":     prog.Description,
		"start_date":      prog.StartDate,
		"end_date":        prog.EndDate,
		"info_url":        prog.InfoURL,
		"wiki_url":        prog.WikiURL,
		"image_url":       prog.Image.ImageURL,
		"image_thumb_url": prog.Image.ThumbnailURL,
		"api_url":         prog.URL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", pbURL, bytes.NewBuffer(jsonData))
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
