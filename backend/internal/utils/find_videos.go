package utils

import (
	"encoding/json"
	"fmt"
	"io"

	// "log"
	"net/http"
)

type LaunchAPIResponse struct {
	Results []Launch `json:"results"`
}

type Launch struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Net     string   `json:"net"`
	URL     string   `json:"url"`
	VidURLs []VidURL `json:"vidURLs"`
}

type VidURL struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	Source    string `json:"source"`
	Publisher string `json:"publisher"`
}

func ListLaunchesWithVidURLs() error {
	apiURL := "https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=15&mode=detailed"
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch launches: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var result LaunchAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	count := 0
	for _, launch := range result.Results {
		if len(launch.VidURLs) > 0 {
			count++
			fmt.Printf("Launch: %s\nID: %s\nDate: %s\nURL: %s\nVideo URLs:\n", launch.Name, launch.ID, launch.Net, launch.URL)
			for _, v := range launch.VidURLs {
				fmt.Printf("  - %s (%s) [%s] %s\n", v.Title, v.Publisher, v.Source, v.URL)
			}
			fmt.Println("--------------------------------------------------")
		}
	}
	fmt.Printf("Total launches with video URLs: %d\n", count)
	return nil
}
