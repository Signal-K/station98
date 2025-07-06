package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type SpaceDevsAstronaut struct {
	Name string `json:"name"`
	Role struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	Priority int `json:"id"`
	Status   struct {
		Name string `json:"name"`
	} `json:"status"`
	InSpace     bool    `json:"in_space"`
	EvaTime     string  `json:"eva_time"`
	TimeInSpace string  `json:"time_in_space"`
	DateOfBirth *string `json:"date_of_birth"`
	DateOfDeath *string `json:"date_of_death"`
	Nationality []struct {
		ID                      int    `json:"id"`
		Name                    string `json:"name"`
		Alpha2Code              string `json:"alpha_2_code"`
		Alpha3Code              string `json:"alpha_3_code"`
		NationalityName         string `json:"nationality_name"`
		NationalityNameComposed string `json:"nationality_name_composed"`
	} `json:"nationality"`
	FirstFlight     string `json:"first_flight"`
	LastFlight      string `json:"last_flight"`
	FlightsCount    int    `json:"flights_count"`
	LandingsCount   int    `json:"landings_count"`
	SpacewalksCount int    `json:"spacewalks_count"`
	IsHuman         bool   `json:"is_human"`
	Agency          struct {
		Name string `json:"name"`
		Type struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
		ID int `json:"id"`
	} `json:"agency"`
	Bio          string `json:"bio"`
	WikipediaURL string `json:"wiki"`
}

type SpaceDevsResponseAstronaut struct {
	Results []SpaceDevsAstronaut `json:"results"`
}

// SyncAstronauts fetches astronaut data from SpaceDevs API and syncs it to Pocketbase
func SyncAstronauts() error {
	fmt.Println("🧑🏼‍🚀 Syncing astronaut candidates....")

	resp, err := http.Get("https://ll.thespacedevs.com/2.3.0/astronauts/?limit=50000&mode=detailed&format=json&ordering=-date_of_birth")
	if err != nil {
		return fmt.Errorf("failed to fetch astronauts: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data SpaceDevsResponseAstronaut
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Helper function to check if astronaut has earthling nationality
	isEarthling := func(nationalities []struct {
		ID                      int    `json:"id"`
		Name                    string `json:"name"`
		Alpha2Code              string `json:"alpha_2_code"`
		Alpha3Code              string `json:"alpha_3_code"`
		NationalityName         string `json:"nationality_name"`
		NationalityNameComposed string `json:"nationality_name_composed"`
	}) bool {
		for _, nat := range nationalities {
			if strings.ToLower(nat.NationalityName) == "earthling" ||
				strings.ToLower(nat.Name) == "earthling" {
				return true
			}
		}
		return false
	}

	var errors []error
	filtered := 0
	for _, astro := range data.Results {
		// Skip astronauts with earthling nationality
		if isEarthling(astro.Nationality) {
			log.Printf("🌍 Skipping earthling: %s", astro.Name)
			filtered++
			continue
		}

		if err := createAstronautInPocketbase(astro); err != nil {
			log.Printf("❌ Error inserting %s: %v", astro.Name, err)
			errors = append(errors, fmt.Errorf("failed to insert %s: %w", astro.Name, err))
		} else {
			log.Printf("✅ Inserted astronaut: %s", astro.Name)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during sync: %v", len(errors), errors[0])
	}

	processed := len(data.Results) - filtered
	fmt.Printf("✅ Successfully processed %d astronauts (filtered out %d earthlings from %d total)\n", processed, filtered, len(data.Results))
	return nil
}

func createAstronautInPocketbase(astro SpaceDevsAstronaut) error {
	pbURL := "http://127.0.0.1:8080/api/collections/astronauts/records"

	// Helper function to safely get string from pointer
	getStringPtr := func(s *string) string {
		if s != nil {
			return *s
		}
		return ""
	}

	// Map API status values to PocketBase allowed values
	mapStatus := func(apiStatus string) string {
		switch strings.ToLower(apiStatus) {
		case "active", "in-training", "occasional spaceflight":
			return "Active"
		case "retired", "deceased":
			return "Retired"
		default:
			// Default to Active for unknown statuses
			return "Active"
		}
	}

	// Helper function to get nationality string
	getNationality := func(nationalities []struct {
		ID                      int    `json:"id"`
		Name                    string `json:"name"`
		Alpha2Code              string `json:"alpha_2_code"`
		Alpha3Code              string `json:"alpha_3_code"`
		NationalityName         string `json:"nationality_name"`
		NationalityNameComposed string `json:"nationality_name_composed"`
	}) string {
		if len(nationalities) > 0 {
			return nationalities[0].NationalityName
		}
		return "Unknown"
	}

	payload := map[string]interface{}{
		"name":             astro.Name,
		"role":             astro.Role.Name,
		"priority":         astro.Priority,
		"status":           mapStatus(astro.Status.Name),
		"in_space":         astro.InSpace,
		"eva_time_total":   astro.EvaTime,
		"space_time_total": astro.TimeInSpace,
		"dob":              getStringPtr(astro.DateOfBirth),
		"nationality":      getNationality(astro.Nationality),
		"first_flight":     astro.FirstFlight,
		"last_flight":      astro.LastFlight,
		"flights_count":    astro.FlightsCount,
		"landings_count":   astro.LandingsCount,
		"spacewalks_count": astro.SpacewalksCount,
		"is_human":         astro.IsHuman,
		"date_of_death":    getStringPtr(astro.DateOfDeath),
		"bio":              astro.Bio,
		"wikipedia_url":    astro.WikipediaURL,
		"agency_type":      astro.Agency.Type.Name,
		// // Assumes 'agency' is a relation field using agency ID or name
		// // Replace with correct Pocketbase ID if needed
		// "agency": astro.Agency.ID,
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
