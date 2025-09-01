package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type LandingAPIResponse struct {
	Results []Landing `json:"results"`
}

type LandingType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type Location struct {
	ID                int    `json:"id"`
	URL               string `json:"url"`
	Name              string `json:"name"`
	CountryCode       string `json:"country_code"`
	Description       string `json:"description"`
	MapImage          string `json:"map_image"`
	TimezoneName      string `json:"timezone_name"`
	TotalLaunchCount  int    `json:"total_launch_count"`
	TotalLandingCount int    `json:"total_landing_count"`
}

type LandingLocation struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Abbrev             string   `json:"abbrev"`
	Description        string   `json:"description"`
	Location           Location `json:"location"`
	SuccessfulLandings int      `json:"successful_landings"`
}

type Launcher struct {
	ID                 int    `json:"id"`
	URL                string `json:"url"`
	Details            string `json:"details"`
	FlightProven       bool   `json:"flight_proven"`
	SerialNumber       string `json:"serial_number"`
	Status             string `json:"status"`
	ImageURL           string `json:"image_url"`
	SuccessfulLandings int    `json:"successful_landings"`
	AttemptedLandings  int    `json:"attempted_landings"`
	Flights            int    `json:"flights"`
	LastLaunchDate     string `json:"last_launch_date"`
	FirstLaunchDate    string `json:"first_launch_date"`
}

type Agency struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	CountryCode  string `json:"country_code"`
	Abbrev       string `json:"abbrev"`
	Description  string `json:"description"`
	FoundingYear string `json:"founding_year"`
	LogoURL      string `json:"logo_url"`
	ImageURL     string `json:"image_url"`
}

type LaunchStatus struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type RocketConfig struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Family   string `json:"family"`
	FullName string `json:"full_name"`
	Variant  string `json:"variant"`
}

type Rocket struct {
	ID            int          `json:"id"`
	Configuration RocketConfig `json:"configuration"`
}

type Orbit struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

type Mission struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	LaunchDesignator string   `json:"launch_designator"`
	Type             string   `json:"type"`
	Orbit            Orbit    `json:"orbit"`
	Agencies         []Agency `json:"agencies"`
}

type Pad struct {
	ID                        int      `json:"id"`
	URL                       string   `json:"url"`
	AgencyID                  int      `json:"agency_id"`
	Name                      string   `json:"name"`
	Description               string   `json:"description"`
	InfoURL                   string   `json:"info_url"`
	WikiURL                   string   `json:"wiki_url"`
	MapURL                    string   `json:"map_url"`
	Latitude                  string   `json:"latitude"`
	Longitude                 string   `json:"longitude"`
	Location                  Location `json:"location"`
	CountryCode               string   `json:"country_code"`
	MapImage                  string   `json:"map_image"`
	TotalLaunchCount          int      `json:"total_launch_count"`
	OrbitalLaunchAttemptCount int      `json:"orbital_launch_attempt_count"`
}

type PreviousFlight struct {
	ID                    string       `json:"id"`
	URL                   string       `json:"url"`
	Slug                  string       `json:"slug"`
	Name                  string       `json:"name"`
	Status                LaunchStatus `json:"status"`
	LastUpdated           string       `json:"last_updated"`
	Net                   string       `json:"net"`
	LaunchServiceProvider Agency       `json:"launch_service_provider"`
	Rocket                Rocket       `json:"rocket"`
	Mission               Mission      `json:"mission"`
	Pad                   Pad          `json:"pad"`
	WebcastLive           bool         `json:"webcast_live"`
	Image                 string       `json:"image"`
}

type FirstStage struct {
	ID                   int             `json:"id"`
	Type                 string          `json:"type"`
	Reused               bool            `json:"reused"`
	LauncherFlightNumber int             `json:"launcher_flight_number"`
	Launcher             Launcher        `json:"launcher"`
	PreviousFlightDate   string          `json:"previous_flight_date"`
	TurnAroundTimeDays   int             `json:"turn_around_time_days"`
	PreviousFlight       *PreviousFlight `json:"previous_flight"`
}

type Landing struct {
	ID                int             `json:"id"`
	URL               string          `json:"url"`
	Attempt           bool            `json:"attempt"`
	Success           *bool           `json:"success"`
	Description       string          `json:"description"`
	DownrangeDistance float64         `json:"downrange_distance"`
	LandingType       LandingType     `json:"landing_type"`
	LandingLocation   LandingLocation `json:"landing_location"`
	FirstStage        FirstStage      `json:"firststage"`
	SpacecraftFlight  interface{}     `json:"spacecraftflight"`
}

func main() {
	fmt.Println("🚀 Fetching most recent landing from Space Devs API...")

	// Hit the Space Devs API for landings (order by ID descending to get most recent)
	url := "https://ll.thespacedevs.com/2.2.0/landings/?limit=1&ordering=-id&mode=detailed&format=json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch landings: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result LandingAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(result.Results) == 0 {
		fmt.Println("❌ No landings found")
		return
	}

	// Display the first (most recent) landing
	landing := result.Results[0]

	fmt.Printf("✅ Most recent landing found!\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════════════\n\n")

	// Basic Landing Info
	fmt.Printf("🎆 LANDING EVENT\n")
	fmt.Printf("├─ ID: %d\n", landing.ID)
	fmt.Printf("├─ URL: %s\n", landing.URL)
	fmt.Printf("├─ Attempt: %t\n", landing.Attempt)
	successStatus := "Unknown"
	if landing.Success != nil {
		if *landing.Success {
			successStatus = "✅ Success"
		} else {
			successStatus = "❌ Failed"
		}
	}
	fmt.Printf("├─ Success: %s\n", successStatus)
	fmt.Printf("├─ Description: %s\n", landing.Description)
	if landing.DownrangeDistance > 0 {
		fmt.Printf("└─ Downrange Distance: %.1f km\n", landing.DownrangeDistance)
	} else {
		fmt.Printf("└─ Downrange Distance: Not available\n")
	}

	fmt.Printf("\n🏁 LANDING TYPE & LOCATION\n")
	fmt.Printf("├─ Type: %s (%s)\n", landing.LandingType.Name, landing.LandingType.Abbrev)
	fmt.Printf("├─ Type Description: %s\n", landing.LandingType.Description)
	fmt.Printf("├─ Location: %s (%s)\n", landing.LandingLocation.Name, landing.LandingLocation.Abbrev)
	fmt.Printf("├─ Location Description: %s\n", landing.LandingLocation.Description)
	fmt.Printf("├─ Successful Landings at Location: %d\n", landing.LandingLocation.SuccessfulLandings)
	fmt.Printf("├─ Geographic Location: %s (%s)\n", landing.LandingLocation.Location.Name, landing.LandingLocation.Location.CountryCode)
	fmt.Printf("├─ Total Launches at Geographic Location: %d\n", landing.LandingLocation.Location.TotalLaunchCount)
	fmt.Printf("└─ Total Landings at Geographic Location: %d\n", landing.LandingLocation.Location.TotalLandingCount)

	fmt.Printf("\n🚀 FIRST STAGE & LAUNCHER\n")
	fmt.Printf("├─ Stage ID: %d (Type: %s)\n", landing.FirstStage.ID, landing.FirstStage.Type)
	fmt.Printf("├─ Reused: %t\n", landing.FirstStage.Reused)
	fmt.Printf("├─ Flight Number for this Launcher: %d\n", landing.FirstStage.LauncherFlightNumber)
	fmt.Printf("├─ Turn Around Time: %d days\n", landing.FirstStage.TurnAroundTimeDays)
	fmt.Printf("├─ Previous Flight Date: %s\n", formatDate(landing.FirstStage.PreviousFlightDate))
	fmt.Printf("├─ Launcher ID: %d\n", landing.FirstStage.Launcher.ID)
	fmt.Printf("├─ Serial Number: %s\n", landing.FirstStage.Launcher.SerialNumber)
	fmt.Printf("├─ Status: %s\n", strings.Title(landing.FirstStage.Launcher.Status))
	fmt.Printf("├─ Flight Proven: %t\n", landing.FirstStage.Launcher.FlightProven)
	fmt.Printf("├─ Details: %s\n", landing.FirstStage.Launcher.Details)
	fmt.Printf("├─ Total Flights: %d\n", landing.FirstStage.Launcher.Flights)
	fmt.Printf("├─ Successful Landings: %d/%d\n", landing.FirstStage.Launcher.SuccessfulLandings, landing.FirstStage.Launcher.AttemptedLandings)
	fmt.Printf("├─ First Launch Date: %s\n", formatDate(landing.FirstStage.Launcher.FirstLaunchDate))
	fmt.Printf("└─ Last Launch Date: %s\n", formatDate(landing.FirstStage.Launcher.LastLaunchDate))

	// Previous Flight Information
	if landing.FirstStage.PreviousFlight != nil {
		pf := landing.FirstStage.PreviousFlight
		fmt.Printf("\n🏁 PREVIOUS LAUNCH\n")
		fmt.Printf("├─ Name: %s\n", pf.Name)
		fmt.Printf("├─ Slug: %s\n", pf.Slug)
		fmt.Printf("├─ Status: %s (%s)\n", pf.Status.Name, pf.Status.Abbrev)
		fmt.Printf("├─ Launch Date: %s\n", formatDate(pf.Net))
		fmt.Printf("├─ Last Updated: %s\n", formatDate(pf.LastUpdated))
		fmt.Printf("├─ Webcast Live: %t\n", pf.WebcastLive)
		fmt.Printf("└─ Launch URL: %s\n", pf.URL)

		// AGENCY/PROVIDER (Highlighted)
		fmt.Printf("\n🏢 AGENCY/PROVIDER\n")
		fmt.Printf("├─ Name: %s\n", pf.LaunchServiceProvider.Name)
		fmt.Printf("├─ Type: %s\n", pf.LaunchServiceProvider.Type)
		if pf.LaunchServiceProvider.CountryCode != "" {
			fmt.Printf("├─ Country: %s\n", pf.LaunchServiceProvider.CountryCode)
		}
		if pf.LaunchServiceProvider.Abbrev != "" {
			fmt.Printf("├─ Abbreviation: %s\n", pf.LaunchServiceProvider.Abbrev)
		}
		if pf.LaunchServiceProvider.FoundingYear != "" {
			fmt.Printf("├─ Founded: %s\n", pf.LaunchServiceProvider.FoundingYear)
		}
		fmt.Printf("├─ Description: %s\n", pf.LaunchServiceProvider.Description)
		fmt.Printf("└─ URL: %s\n", pf.LaunchServiceProvider.URL)

		// ROCKET (Highlighted)
		fmt.Printf("\n🚀 ROCKET\n")
		fmt.Printf("├─ Configuration: %s\n", pf.Rocket.Configuration.FullName)
		fmt.Printf("├─ Family: %s\n", pf.Rocket.Configuration.Family)
		fmt.Printf("├─ Variant: %s\n", pf.Rocket.Configuration.Variant)
		fmt.Printf("└─ Config URL: %s\n", pf.Rocket.Configuration.URL)

		// MISSION (Highlighted)
		fmt.Printf("\n🌌 MISSION\n")
		fmt.Printf("├─ Name: %s\n", pf.Mission.Name)
		fmt.Printf("├─ Type: %s\n", pf.Mission.Type)
		fmt.Printf("├─ Orbit: %s (%s)\n", pf.Mission.Orbit.Name, pf.Mission.Orbit.Abbrev)
		if pf.Mission.LaunchDesignator != "" {
			fmt.Printf("├─ Launch Designator: %s\n", pf.Mission.LaunchDesignator)
		}
		fmt.Printf("├─ Description: %s\n", pf.Mission.Description)
		if len(pf.Mission.Agencies) > 0 {
			fmt.Printf("└─ Mission Agencies:\n")
			for i, agency := range pf.Mission.Agencies {
				prefix := "   ├─"
				if i == len(pf.Mission.Agencies)-1 {
					prefix = "   └─"
				}
				fmt.Printf("%s %s (%s) - %s\n", prefix, agency.Name, agency.Abbrev, agency.CountryCode)
			}
		} else {
			fmt.Printf("└─ Mission Agencies: None specified\n")
		}

		// PAD/LAUNCHPAD (Highlighted)
		fmt.Printf("\n🗺 PAD/LAUNCHPAD\n")
		fmt.Printf("├─ Name: %s\n", pf.Pad.Name)
		fmt.Printf("├─ Location: %s (%s)\n", pf.Pad.Location.Name, pf.Pad.Location.CountryCode)
		fmt.Printf("├─ Coordinates: %s, %s\n", pf.Pad.Latitude, pf.Pad.Longitude)
		fmt.Printf("├─ Agency ID: %d\n", pf.Pad.AgencyID)
		if pf.Pad.Description != "" {
			fmt.Printf("├─ Description: %s\n", pf.Pad.Description)
		}
		fmt.Printf("├─ Total Launch Count: %d\n", pf.Pad.TotalLaunchCount)
		fmt.Printf("├─ Orbital Launch Attempts: %d\n", pf.Pad.OrbitalLaunchAttemptCount)
		if pf.Pad.WikiURL != "" {
			fmt.Printf("├─ Wiki: %s\n", pf.Pad.WikiURL)
		}
		if pf.Pad.MapURL != "" {
			fmt.Printf("├─ Map: %s\n", pf.Pad.MapURL)
		}
		fmt.Printf("└─ Pad URL: %s\n", pf.Pad.URL)
	}

	fmt.Printf("\n═══════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("🎯 Landing data retrieved successfully!\n")
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return "Not available"
	}
	return dateStr
}
