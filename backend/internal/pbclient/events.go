package pbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Event struct {
	Title       string `json:"title"`
	Datetime    string `json:"datetime"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	SourceURL   string `json:"source_url"`
	Description string `json:"description"`
	SpacedevsID string `json:"spacedevs_id"` // string, not a relation
	ProviderID  string `json:"provider"`     // relation to launch_providers
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

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to insert event (%s): %s", e.Title, resp.Status)
	}
	return nil
}
