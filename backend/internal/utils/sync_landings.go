package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/signal-k/notifs/internal/pbclient"
	"github.com/signal-k/notifs/internal/sync"
)

// LandingAPIResponse matches the Space Devs API response for landings
// Only the fields needed for syncing are included here
// Extend as needed for your use case

type LandingAPIResponse struct {
	Results []Landing `json:"results"`
}

type Landing struct {
	ID               int    `json:"id"`
	Attempt          bool   `json:"attempt"`
	Success          bool   `json:"success"`
	LandingType      string `json:"landing_type__name"`
	LandingLocation  string `json:"landing_location__name"`
	LandingDate      string `json:"date"`
	FirstStageLaunch int    `json:"firststage_launch__id"`
	LauncherID       int    `json:"launcher__id"`
	RocketName       string `json:"rocket__launch__name"`
	Agency           sync.Agency
	// Add more fields as needed
}

// SyncMostRecentLanding fetches the most recent landing and returns its data
func SyncMostRecentLanding(client *pbclient.Client) (*Landing, error) {
	url := "https://ll.thespacedevs.com/2.2.0/landings/?limit=1&ordering=-date&mode=detailed&format=json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, sync.NewSyncError(resp.StatusCode, string(body))
	}

	var result LandingAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, nil // No landings found
	}

	landing := result.Results[0]
	log.Printf("Most recent landing: %+v", landing)
	return &landing, nil
}
