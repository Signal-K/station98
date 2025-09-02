package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExpeditionAPIResponse struct {
	Count   int          `json:"count"`
	Next    string       `json:"next"`
	Prev    string       `json:"previous"`
	Results []Expedition `json:"results"`
}

// Expedition struct (detailed mode)
type Expedition struct {
	ID             int          `json:"id"`
	URL            string       `json:"url"`
	Name           string       `json:"name"`
	Start          string       `json:"start"`
	End            string       `json:"end"`
	SpaceStation   SpaceStation `json:"spacestation"`
	Crew           []CrewMember `json:"crew"`
	MissionPatches []Patch      `json:"mission_patches"`
}

type NamedObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SpaceStation struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Role struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Priority int    `json:"priority"`
}

type CrewMember struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Astronaut Astronaut `json:"astronaut"`
}

type Nationality struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Alpha2Code             string `json:"alpha_2_code"`
	Alpha3Code             string `json:"alpha_3_code"`
	NationalityName        string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Agency struct {
	ResponseMode string     `json:"response_mode"`
	ID           int        `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	Abbrev       string     `json:"abbrev"`
	Type         NamedObject `json:"type"`
}

type Astronaut struct {
	ID           int           `json:"id"`
	URL          string        `json:"url"`
	Name         string        `json:"name"`
	Status       Status        `json:"status"`
	Age          int           `json:"age"`
	InSpace      bool          `json:"in_space"`
	TimeInSpace  string        `json:"time_in_space"`
	EVATime      string        `json:"eva_time"`
	DateOfBirth  string        `json:"date_of_birth"`
	DateOfDeath  *string       `json:"date_of_death"`
	Nationality  []Nationality `json:"nationality"`
	Agency       Agency        `json:"agency"`
	FirstFlight  string        `json:"first_flight"`
	LastFlight   string        `json:"last_flight"`
	Wiki         *string       `json:"wiki"`
	Bio          string        `json:"bio"`
}

type Patch struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	ImageURL string `json:"image_url"`
}

func main() {
	expeditions, err := GetRecentExpeditions()
	if err != nil {
		fmt.Println("❌ Error fetching expeditions:", err)
		return
	}

	fmt.Println("🧭 Most Recent Space Expeditions:\n")

	for _, e := range expeditions {
		fmt.Println("════════════════════════════════════════════════════")
		fmt.Printf("🚀 Name: %s\n", e.Name)
		fmt.Printf("📅 Start: %s\n", e.Start)
		fmt.Printf("📅 End: %s\n", e.End)
		fmt.Printf("🏢 Space Station: %s\n", e.SpaceStation.Name)
		if len(e.MissionPatches) > 0 {
			fmt.Printf("🎖️ Patch: %s (%s)\n", e.MissionPatches[0].Name, e.MissionPatches[0].ImageURL)
		}
		fmt.Println("👩‍🚀 Crew:")
		for _, member := range e.Crew {
			astro := member.Astronaut
			nationalityStr := "Unknown"
			if len(astro.Nationality) > 0 {
				nationalityStr = astro.Nationality[0].NationalityName
			}
			fmt.Printf("   - %s (%s, %s, Age %d)\n", astro.Name, member.Role.Role, nationalityStr, astro.Age)
		}
	}
}

func GetRecentExpeditions() ([]Expedition, error) {
	url := "https://ll.thespacedevs.com/2.3.0/expeditions/?limit=5&ordering=-start&mode=detailed&format=json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed ExpeditionAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("Error decoding: %w", err)
	}

	return parsed.Results, nil
}
