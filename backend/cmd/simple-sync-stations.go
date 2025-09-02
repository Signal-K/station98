package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Top-level API response
type SpaceStationAPIResponse struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Prev    string         `json:"previous"`
	Results []SpaceStation `json:"results"`
}

func main() {
	stations, err := GetAllNewestSpaceStations()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, s := range stations {
		fmt.Println("════════════════════════════════════════════════════")
		fmt.Printf("🛰️ Name: %s\n", s.Name)
		fmt.Printf("📅 Founded: %s\n", s.Founded)
		fmt.Printf("📝 Description: %s\n", s.Description)
		fmt.Printf("🧑‍🚀 Onboard Crew: %d\n", s.OnboardCrew)
		fmt.Printf("🚢 Docked Vehicles: %d\n", s.DockedVehicles)
		fmt.Printf("📍 Orbit: %s\n", s.Orbit)
		fmt.Printf("🏗️ Type: %s\n", s.Type.Name)
		fmt.Printf("📊 Mass: %.0f kg\n", s.Mass)
		fmt.Printf("📐 Width: %.0f m, Height: %.0f m, Volume: %.0f m³\n", s.Width, s.Height, s.Volume)
		fmt.Printf("🔗 URL: %s\n", s.URL)

		// Owners
		fmt.Println("🏛️ Owners:")
		for _, o := range s.Owners {
			fmt.Printf("   - %s (%s)\n", o.Name, o.Type.Name)
		}

		// Expeditions
		fmt.Println("🚀 Expeditions:")
		for _, e := range s.Expeditions {
			fmt.Printf("   - %s (%s → %s)\n", e.Name, e.Start, e.End)
		}

		// Docking locations
		fmt.Println("🛰️ Docking Locations:")
		for _, d := range s.DockingLocation {
			fmt.Printf("   - %s\n", d.Name)
			if d.CurrentlyDocked != nil {
				fmt.Printf("     🚀 Docked Vehicle ID: %d\n", d.CurrentlyDocked.ID)
			}
		}
	}
}

// Deep struct with full detail
type SpaceStation struct {
	ID              int               `json:"id"`
	URL             string            `json:"url"`
	Name            string            `json:"name"`
	Image           Image             `json:"image"`
	Status          NamedObject       `json:"status"`
	Founded         string            `json:"founded"`
	Deorbited       string            `json:"deorbited"`
	Description     string            `json:"description"`
	Orbit           string            `json:"orbit"`
	Type            NamedObject       `json:"type"`
	Owners          []Agency          `json:"owners"`
	Expeditions     []Expedition      `json:"active_expeditions"`
	Height          float64           `json:"height"`
	Width           float64           `json:"width"`
	Mass            float64           `json:"mass"`
	Volume          float64           `json:"volume"`
	OnboardCrew     int               `json:"onboard_crew"`
	DockedVehicles  int               `json:"docked_vehicles"`
	DockingLocation []DockingLocation `json:"docking_location"`
}

type NamedObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Image struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Credit       string    `json:"credit"`
	License      License   `json:"license"`
	SingleUse    bool      `json:"single_use"`
	Variants     []Variant `json:"variants"`
}

type Variant struct {
	ID       int         `json:"id"`
	Type     NamedObject `json:"type"`
	ImageURL string      `json:"image_url"`
}

type License struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Link     string `json:"link"`
}

type Agency struct {
	ID         int         `json:"id"`
	URL        string      `json:"url"`
	Name       string      `json:"name"`
	Abbrev     string      `json:"abbrev"`
	Type       NamedObject `json:"type"`
	Featured   bool        `json:"featured"`
	Country    []Country   `json:"country"`
	Logo       Image       `json:"logo"`
	Image      Image       `json:"image"`
	SocialLogo Image       `json:"social_logo"`
}

type Country struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	Alpha2Code              string `json:"alpha_2_code"`
	Alpha3Code              string `json:"alpha_3_code"`
	NationalityName         string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type Expedition struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type DockingLocation struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	CurrentlyDocked *DockedVehicle `json:"currently_docked"`
}

type DockedVehicle struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Docking   string `json:"docking"`
	Departure string `json:"departure"`
}

// Fetch the 2 most recently founded space stations
func GetAllNewestSpaceStations() ([]SpaceStation, error) {
	url := "https://ll.thespacedevs.com/2.3.0/space_stations/?limit=2&ordering=-founded&mode=detailed"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed SpaceStationAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	return parsed.Results, nil
}
