package pbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Update struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type Link struct {
	Priority string `json:"priority"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type TimelineEntry struct {
	Time  int    `json:"time"`
	Event string `json:"event"`
}

type Event struct {
	Title             string          `json:"title"`
	Datetime          string          `json:"datetime"`
	Location          string          `json:"location"`
	Type              string          `json:"type"`
	SourceURL         string          `json:"source_url"`
	Description       string          `json:"description"`
	SpacedevsID       string          `json:"spacedevs_id"`
	ProviderID        string          `json:"provider"`
	StatusAbbrev      string          `json:"status_abbrev"`
	StatusDescription string          `json:"status_description"`
	WeatherConcerns   string          `json:"weather_concerns"`
	Updates           []Update        `json:"updates"`
	RocketID          string          `json:"rocket_id"`
	MissionID         string          `json:"mission_id"`
	PadID             string          `json:"pad_id"`
	InfoURLs          []Link          `json:"info_urls"`
	VideoURLs         []Link          `json:"vid_urls"`
	WebcastLive       bool            `json:"webcast_live"`
	Timeline          []TimelineEntry `json:"timeline"`
}

func (c *Client) CreateEvent(e Event) error {
	data, _ := json.Marshal(e)

	req, _ := http.NewRequest("POST", c.BaseURL+"/api/collections/events/records", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	existing, err := c.FindRecordByField("events", "title", e.Title)
	if err == nil && existing != nil {
		return fmt.Errorf("event %s already exists", e.Title)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to insert event (%s): %s", e.Title, resp.Status)
	}
	return nil
}
