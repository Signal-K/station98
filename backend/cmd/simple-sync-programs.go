package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ProgramResponse struct {
	Count   int       `json:"count"`
	Next    string    `json:"next"`
	Results []Program `json:"results"`
}

type Program struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Image       Image     `json:"image"`
	InfoURL     string    `json:"info_url"`
	WikiURL     string    `json:"wiki_url"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Type        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
}

type Image struct {
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Name         string `json:"name"`
	Credit       string `json:"credit"`
}

func fetchLatestPrograms() error {
	url := "https://ll.thespacedevs.com/2.2.0/program/?limit=5&ordering=-start_date&mode=normal"

	fmt.Println("🚀 Fetching latest space programs...")
	fmt.Printf("🌐 API URL: %s\n\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch programs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var programs ProgramResponse
	if err := json.NewDecoder(resp.Body).Decode(&programs); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	fmt.Printf("✅ Successfully fetched %d programs (Total available: %d)\n\n", len(programs.Results), programs.Count)
	fmt.Println("🚀 Latest Space Programs:")
	fmt.Println("════════════════════════════════════════════════")

	for i, p := range programs.Results {
		fmt.Printf("\n💫 Program #%d\n", i+1)
		fmt.Printf("📛 Name: %s (ID: %d)\n", p.Name, p.ID)
		fmt.Printf("🗺️  Duration: %s to %s\n", p.StartDate.Format("2006-01-02"), p.EndDate.Format("2006-01-02"))
		fmt.Printf("🧭 Type: %s\n", p.Type.Name)
		
		if p.Description != "" {
			// Truncate description if too long
			desc := p.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			fmt.Printf("📜 Description: %s\n", desc)
		}
		
		if p.InfoURL != "" {
			fmt.Printf("🔗 Info URL: %s\n", p.InfoURL)
		}
		
		if p.WikiURL != "" {
			fmt.Printf("🌐 Wiki: %s\n", p.WikiURL)
		}
		
		if p.Image.ImageURL != "" {
			fmt.Printf("🖼️  Image: %s\n", p.Image.ImageURL)
			if p.Image.Credit != "" {
				fmt.Printf("📷 Image Credit: %s\n", p.Image.Credit)
			}
		}
		
		fmt.Printf("🌐 API URL: %s\n", p.URL)
		
		if i < len(programs.Results)-1 {
			fmt.Println("────────────────────────────────────────────────")
		}
	}

	fmt.Println("\n════════════════════════════════════════════════")

	if programs.Next != "" {
		fmt.Printf("➡️  Next page: %s\n", programs.Next)
	}

	return nil
}

func main() {
	if err := fetchLatestPrograms(); err != nil {
		fmt.Println("Error:", err)
	}
}
