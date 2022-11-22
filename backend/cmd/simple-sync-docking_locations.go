package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DockingLocationResponse struct {
	Count   int               `json:"count"`
	Results []DockingLocation `json:"results"`
}

type DockingLocation struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	SpaceStation SpaceStation `json:"spacestation"`
	Spacecraft   *Spacecraft  `json:"spacecraft"`
	Payload      *Payload     `json:"payload"`
}

type SpaceStation struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	URL   string      `json:"url"`
	Image ImageDetail `json:"image"`
}

type Spacecraft struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ImageDetail struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	ImageURL     string  `json:"image_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Credit       string  `json:"credit"`
	License      License `json:"license"`
	SingleUse    bool    `json:"single_use"`
}

type License struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Link     string `json:"link"`
}

func main() {
	fmt.Println("🛰️ Fetching most recent docking location...")

	url := "https://ll.thespacedevs.com/2.3.0/config/docking_locations/?limit=1&format=json"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("❌ Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Unexpected response code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read body: %v\n", err)
		return
	}

	var data DockingLocationResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("❌ JSON error: %v\n", err)
		return
	}

	if len(data.Results) == 0 {
		fmt.Println("⚠️ No docking locations found.")
		return
	}

	loc := data.Results[0]
	fmt.Println("✅ Most Recent Docking Location:")
	fmt.Printf("- Name: %s\n", loc.Name)
	fmt.Printf("- Station: %s\n", loc.SpaceStation.Name)
	fmt.Printf("  • Image: %s\n", loc.SpaceStation.Image.ImageURL)
	fmt.Printf("  • Credit: %s\n", loc.SpaceStation.Image.Credit)
	fmt.Printf("  • License: %s (%s)\n", loc.SpaceStation.Image.License.Name, loc.SpaceStation.Image.License.Link)

	if loc.Spacecraft != nil {
		fmt.Printf("- Spacecraft: %s\n", loc.Spacecraft.Name)
	} else {
		fmt.Println("- Spacecraft: None")
	}

	if loc.Payload != nil {
		fmt.Printf("- Payload: %s\n", loc.Payload.Name)
	} else {
		fmt.Println("- Payload: None")
	}
}
