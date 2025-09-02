package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AstronautAPIResponse struct {
	Count   int         `json:"count"`
	Next    string      `json:"next"`
	Prev    string      `json:"previous"`
	Results []Astronaut `json:"results"`
}

type Nationality struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Alpha2Code             string `json:"alpha_2_code"`
	Alpha3Code             string `json:"alpha_3_code"`
	NationalityName        string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type AgencyDetail struct {
	ResponseMode string      `json:"response_mode"`
	ID           int         `json:"id"`
	URL          string      `json:"url"`
	Name         string      `json:"name"`
	Abbrev       string      `json:"abbrev"`
	Type         NamedObject `json:"type"`
}

type Astronaut struct {
	ID              int           `json:"id"`
	URL             string        `json:"url"`
	Name            string        `json:"name"`
	Status          NamedObject   `json:"status"`
	Agency          AgencyDetail  `json:"agency"`
	Type            NamedObject   `json:"type"`
	InSpace         bool          `json:"in_space"`
	TimeInSpace     string        `json:"time_in_space"`
	EvaTime         string        `json:"eva_time"`
	Age             *int          `json:"age"`
	DateOfBirth     *string       `json:"date_of_birth"`
	DateOfDeath     *string       `json:"date_of_death"`
	Nationality     []Nationality `json:"nationality"`
	Bio             string        `json:"bio"`
	Wiki            *string       `json:"wiki"`
	LastFlight      *string       `json:"last_flight"`
	FirstFlight     *string       `json:"first_flight"`
	FlightsCount    int           `json:"flights_count"`
	LandingsCount   int           `json:"landings_count"`
	SpacewalksCount int           `json:"spacewalks_count"`
}

type NamedObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	astronauts, err := GetRecentAstronauts()
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("🧑‍🚀 Most Recently Born Astronauts (Newest)")
	for _, a := range astronauts {
		fmt.Println("═══════════════════════════════════════════")
		fmt.Printf("👤 Name: %s\n", a.Name)
		
		// Handle nationality array
		nationalityStr := "Unknown"
		if len(a.Nationality) > 0 {
			nationalities := make([]string, len(a.Nationality))
			for i, nat := range a.Nationality {
				nationalities[i] = nat.NationalityName
			}
			nationalityStr = fmt.Sprintf("%v", nationalities)
			if len(a.Nationality) == 1 {
				nationalityStr = a.Nationality[0].NationalityName
			}
		}
		fmt.Printf("🗺️ Nationality: %s\n", nationalityStr)
		
		// Handle nullable fields
		dob := "Unknown"
		if a.DateOfBirth != nil {
			dob = *a.DateOfBirth
		}
		fmt.Printf("🎂 DOB: %s\n", dob)
		
		dod := "N/A"
		if a.DateOfDeath != nil && *a.DateOfDeath != "" {
			dod = *a.DateOfDeath
		}
		fmt.Printf("☠️ DOD: %s\n", dod)
		
		age := "Unknown"
		if a.Age != nil {
			age = fmt.Sprintf("%d", *a.Age)
		}
		fmt.Printf("🎂 Age: %s\n", age)
		
		fmt.Printf("🧠 Status: %s | Type: %s | In Space Now: %t\n", a.Status.Name, a.Type.Name, a.InSpace)
		fmt.Printf("🏢 Agency: %s (%s)\n", a.Agency.Name, a.Agency.Abbrev)
		fmt.Printf("🚀 Flights: %d | Landings: %d | Spacewalks: %d\n", a.FlightsCount, a.LandingsCount, a.SpacewalksCount)
		fmt.Printf("🕒 Time in Space: %s | EVA Time: %s\n", a.TimeInSpace, a.EvaTime)
		
		firstFlight := "N/A"
		if a.FirstFlight != nil && *a.FirstFlight != "" {
			firstFlight = *a.FirstFlight
		}
		lastFlight := "N/A"
		if a.LastFlight != nil && *a.LastFlight != "" {
			lastFlight = *a.LastFlight
		}
		fmt.Printf("🛰️ First Flight: %s | Last Flight: %s\n", firstFlight, lastFlight)
		
		if a.Bio != "" {
			bioPreview := a.Bio
			if len(bioPreview) > 300 {
				bioPreview = bioPreview[:300] + "..."
			}
			fmt.Printf("📜 Bio: %s\n", bioPreview)
		}
		
		wiki := "N/A"
		if a.Wiki != nil && *a.Wiki != "" {
			wiki = *a.Wiki
		}
		fmt.Printf("🌐 Wiki: %s\n", wiki)
	}
}

func GetRecentAstronauts() ([]Astronaut, error) {
	url := "https://ll.thespacedevs.com/2.3.0/astronauts/?limit=3&ordering=-date_of_birth&mode=detailed&format=json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed AstronautAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	return parsed.Results, nil
}
