package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LaunchEE struct {
	Name   string `json:"name"`
	Net    string `json:"net"` // Launch datetime
	Status struct {
		Name string `json:"name"`
	} `json:"status"`
	Mission struct {
		Name string `json:"name"`
	} `json:"mission"`
}

type Expedition struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
	URL   string `json:"url"`
}

type SpaceDevsExpeditionsResponse struct {
	Results []Expedition `json:"results"`
}

type LaunchesResponse struct {
	Results []LaunchEE `json:"results"`
}

// MatchExpeditionLaunches compares expedition start/end with launches on same day.
func MatchExpeditionLaunches() error {
	fmt.Println("📡 Fetching launches and expeditions...")

	expeditions, err := fetchExpeditions()
	if err != nil {
		return err
	}

	launches, err := fetchLaunches()
	if err != nil {
		return err
	}

	fmt.Printf("🔍 Matching %d expeditions against %d launches...\n", len(expeditions), len(launches))

	for _, exp := range expeditions {
		startDate, err := parseDate(exp.Start)
		if err != nil {
			fmt.Printf("⚠️  Invalid start date for expedition %s: %s\n", exp.Name, exp.Start)
			continue
		}
		var endDate time.Time
		if exp.End != "" {
			endDate, _ = parseDate(exp.End)
		}

		matchFound := false

		for _, launch := range launches {
			launchDate, err := parseDate(launch.Net)
			if err != nil {
				fmt.Printf("⚠️  Invalid launch date for '%s': %s\n", launch.Name, launch.Net)
				continue
			}

			// fmt.Printf("🔎 Comparing Expedition '%s' (%s) with Launch '%s' (%s)\n",
			// 	exp.Name, startDate.Format("2006-01-02"), launch.Name, launchDate.Format("2006-01-02"))

			if sameDay(launchDate, startDate) {
				fmt.Printf("🚀 Launch '%s' matches expedition START '%s'\n", launch.Name, exp.Name)
				matchFound = true
			}

			if !endDate.IsZero() && sameDay(launchDate, endDate) {
				fmt.Printf("🛬 Launch '%s' matches expedition END '%s'\n", launch.Name, exp.Name)
				matchFound = true
			}
		}

		if !matchFound {
			fmt.Printf("❌ No launch match found for expedition: %s\n", exp.Name)
		}
	}

	return nil
}

func fetchExpeditions() ([]Expedition, error) {
	url := "https://ll.thespacedevs.com/2.3.0/expeditions/?limit=100&ordering=-start"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch expeditions: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data SpaceDevsExpeditionsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse expeditions JSON: %w", err)
	}
	return data.Results, nil
}

func fetchLaunches() ([]LaunchEE, error) {
	url := "https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=50&mode=detailed&offset="
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch launches: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data LaunchesResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse launches JSON: %w", err)
	}
	return data.Results, nil
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", dateStr)
}
