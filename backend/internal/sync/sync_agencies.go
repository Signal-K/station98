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

type SpaceDevsAgencyResponse struct {
	Results []SpaceDevsAgency `json:"results"`
	Next    string            `json:"next"`
	Count   int               `json:"count"`
}

type SpaceDevsAgency struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Abbrev        string `json:"abbrev"`
	Description   string `json:"description"`
	Administrator string `json:"administrator"`
	FoundingYear  int    `json:"founding_year"`
	Launchers     string `json:"launchers"`
	Spacecraft    string `json:"spacecraft"`
	Featured      bool   `json:"featured"`
	URL           string `json:"url"`

	Type struct {
		Name string `json:"name"`
	} `json:"type"`

	Country []struct {
		Name        string `json:"name"`
		Alpha2      string `json:"alpha_2_code"`
		Nationality string `json:"nationality_name"`
	} `json:"country"`

	Logo struct {
		ImageURL string `json:"image_url"`
	} `json:"logo"`

	SocialLogo struct {
		ImageURL string `json:"image_url"`
	} `json:"social_logo"`

	Image *struct {
		ImageURL string `json:"image_url"`
	} `json:"image"`
}

// SyncAgencies fetches and stores agencies in Pocketbase
func SyncAgencies() error {
	fmt.Println("🏢 Syncing launch agencies...")

	pageURL := "https://ll.thespacedevs.com/2.3.0/agencies/?limit=100"
	client := &http.Client{Timeout: 10 * time.Second}

	var count, skipped int
	for pageURL != "" {
		resp, err := client.Get(pageURL)
		if err != nil {
			return fmt.Errorf("failed to fetch from SpaceDevs: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		var data SpaceDevsAgencyResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		for _, agency := range data.Results {
			if agency.Name == "" {
				skipped++
				continue
			}
			if err := createAgencyInPocketbase(agency); err != nil {
				log.Printf("❌ Failed: %s — %v", agency.Name, err)
			} else {
				log.Printf("✅ Synced agency: %s", agency.Name)
				count++
			}
		}

		pageURL = data.Next
	}

	fmt.Printf("🎯 Completed. Synced %d agencies, skipped %d\n", count, skipped)
	return nil
}

func createAgencyInPocketbase(a SpaceDevsAgency) error {
	pbURL := "http://127.0.0.1:8080/api/collections/agencies/records"

	countryName := ""
	countryCode := ""
	nationality := ""
	if len(a.Country) > 0 {
		countryName = a.Country[0].Name
		countryCode = a.Country[0].Alpha2
		nationality = a.Country[0].Nationality
	}

	payload := map[string]any{
		"api_id":           a.ID,
		"name":             a.Name,
		"abbrev":           a.Abbrev,
		"type_name":        a.Type.Name,
		"description":      a.Description,
		"administrator":    a.Administrator,
		"founding_year":    a.FoundingYear,
		"launchers":        a.Launchers,
		"spacecraft":       a.Spacecraft,
		"featured":         a.Featured,
		"url":              a.URL,
		"country_name":     countryName,
		"country_code":     countryCode,
		"nationality_name": nationality,
		"logo_url":         a.Logo.ImageURL,
		"social_logo_url":  a.SocialLogo.ImageURL,
		"image_url":        "",
	}

	if a.Image != nil {
		payload["image_url"] = a.Image.ImageURL
	}

	body, err := json.Marshal(payload)
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
		return fmt.Errorf("PB error: %d - %s", res.StatusCode, string(body))
	}

	return nil
}
