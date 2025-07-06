package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DockingEventResponse struct {
	Results []DockingEvent `json:"results"`
}

type DockingEvent struct {
	ID        int           `json:"id"`
	URL       string        `json:"url"`
	Docking   string        `json:"docking"`
	Departure string        `json:"departure"`
	Location  Location      `json:"docking_location"`
	Target    PayloadFlight `json:"payload_flight_target"`
	Chaser    PayloadFlight `json:"payload_flight_chaser"`
}

type Location struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Payload *Payload `json:"payload"`
}

type PayloadFlight struct {
	ID      int      `json:"id"`
	Payload *Payload `json:"payload"`
	Launch  *Launch  `json:"launch"`
}

type Payload struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Operator *SpaceDevsAgency `json:"operator"`
	Image    *Image           `json:"image"`
}

type Image struct {
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func SyncDockingEvents() error {
	fmt.Println("🔄 Syncing docking events...")

	url := "https://ll.thespacedevs.com/2.3.0/docking_events/?format=json&limit=1000"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch docking events: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data DockingEventResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	for _, event := range data.Results {
		if err := upsertDockingEvent(event); err != nil {
			log.Printf("❌ Error inserting docking event %d: %v", event.ID, err)
		} else {
			log.Printf("✅ Synced docking event %d", event.ID)
		}
	}

	return nil
}

func upsertDockingEvent(event DockingEvent) error {
	apiURL := "http://127.0.0.1:8080/api/collections/docking_events/records"

	getPayloadInfo := func(p *Payload) (int, string, string, string) {
		if p == nil {
			return 0, "", "", ""
		}
		operator := ""
		if p.Operator != nil {
			operator = p.Operator.Name
		}
		image := ""
		if p.Image != nil {
			image = p.Image.ImageURL
		}
		return p.ID, p.Name, operator, image
	}

	locationPayloadID, locationPayloadName, locationOperator, locationImage := getPayloadInfo(event.Location.Payload)
	chaserPayloadID, chaserPayloadName, chaserOperator, chaserImage := getPayloadInfo(event.Chaser.Payload)
	targetPayloadID, targetPayloadName, targetOperator, targetImage := getPayloadInfo(event.Target.Payload)

	chaserLaunchID, chaserLaunchName := "", ""
	if event.Chaser.Launch != nil {
		chaserLaunchID = event.Chaser.Launch.ID
		chaserLaunchName = event.Chaser.Launch.Name
	}

	targetLaunchID, targetLaunchName := "", ""
	if event.Target.Launch != nil {
		targetLaunchID = event.Target.Launch.ID
		targetLaunchName = event.Target.Launch.Name
	}

	payload := map[string]interface{}{
		"api_id":                event.ID,
		"docking_time":          event.Docking,
		"departure_time":        event.Departure,
		"location_name":         event.Location.Name,
		"location_payload_id":   locationPayloadID,
		"location_payload_name": locationPayloadName,
		"location_operator":     locationOperator,
		"location_image_url":    locationImage,
		"chaser_payload_id":     chaserPayloadID,
		"chaser_payload_name":   chaserPayloadName,
		"chaser_operator":       chaserOperator,
		"chaser_image_url":      chaserImage,
		"chaser_launch_id":      chaserLaunchID,
		"chaser_launch_name":    chaserLaunchName,
		"target_payload_id":     targetPayloadID,
		"target_payload_name":   targetPayloadName,
		"target_operator":       targetOperator,
		"target_image_url":      targetImage,
		"target_launch_id":      targetLaunchID,
		"target_launch_name":    targetLaunchName,
		"is_active":             true,
		"details":               fmt.Sprintf("Docking with %s", locationPayloadName),
		"source_url":            event.URL,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("PB insert error: %s", string(body))
	}

	return nil
}
