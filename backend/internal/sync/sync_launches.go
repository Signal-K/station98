package sync

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/signal-k/notifs/internal/pbclient"
)

type Agency struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Abbrev   string `json:"abbrev"`
	Type     string `json:"type"`
	Country  string `json:"country_code"`
	Desc     string `json:"description"`
	Founded  string `json:"founding_year"`
	LogoURL  string `json:"logo_url"`
	ImageURL string `json:"image_url"`
	WikiURL  string `json:"wiki_url"`
	InfoURL  string `json:"info_url"`
}

type Launch struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Net                   string          `json:"net"`
	URL                   string          `json:"url"`
	Status                Status          `json:"status"`
	LaunchServiceProvider Agency          `json:"launch_service_provider"`
	Mission               *Mission        `json:"mission"`
	Rocket                Rocket          `json:"rocket"`
	Pad                   Pad             `json:"pad"`
	InfoURLs              []Link          `json:"infoURLs"`
	VideoURLs             []Link          `json:"vidURLs"`
	WebcastLive           bool            `json:"webcast_live"`
	Timeline              []TimelineEntry `json:"timeline"`
	Updates               []Update        `json:"updates"`
}

type Status struct {
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type Link struct {
	Priority int    `json:"priority"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type TimelineEntry struct {
	Time  int    `json:"time"` // seconds before/after T0
	Event string `json:"event"`
}

type Update struct {
	ID           int    `json:"id"`
	ProfileImage string `json:"profile_image"`
	Comment      string `json:"comment"`
	InfoURL      string `json:"info_url"`
	CreatedBy    string `json:"created_by"`
	CreatedOn    string `json:"created_on"`
}

type Orbit struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

type Mission struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Orbit       Orbit    `json:"orbit"`
	Agencies    []Agency `json:"agencies"`
}

type Rocket struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	Variant      string `json:"variant"`
	Family       string `json:"family"`
	Reusable     bool   `json:"reusable"`
	Description  string `json:"description"`
	LaunchMass   int    `json:"launch_mass"`
	LEOCapacity  int    `json:"leo_capacity"`
	GTOCapacity  int    `json:"gto_capacity"`
	ImageURL     string `json:"image_url"`
	InfoURL      string `json:"info_url"`
	WikiURL      string `json:"wiki_url"`
	Manufacturer Agency `json:"manufacturer"`
}

type Pad struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	CountryCode      string `json:"country_code"`
	LocationName     string `json:"location_name"`
	MapURL           string `json:"map_url"`
	WikiURL          string `json:"wiki_url"`
	MapImage         string `json:"map_image"`
	LaunchCountYear  int    `json:"launch_count_year"`
	LaunchCountTotal int    `json:"launch_count_total"`
}

type LaunchAPIResponse struct {
	Results []Launch `json:"results"`
}

func getThrottleDelay(body string) time.Duration {
	// Try to extract seconds from the error message
	if idx := strings.Index(body, "Expected available in "); idx != -1 {
		rest := body[idx+21:]
		secStr := ""
		for _, c := range rest {
			if c >= '0' && c <= '9' {
				secStr += string(c)
			} else {
				break
			}
		}
		if secStr != "" {
			secs, err := strconv.Atoi(secStr)
			if err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}
	return 60 * time.Second // fallback
}

func SyncLaunchProvidersAndEvents(client *pbclient.Client) {
	log.Println("📡 Fetching upcoming launches...")

	resp, err := http.Get("https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=50&mode=detailed")
	if err != nil {
		log.Printf("❌ Failed to fetch launches: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		body, _ := io.ReadAll(resp.Body)
		delay := getThrottleDelay(string(body))
		log.Printf("❌ Launches: HTTP 429: %s. Retrying in %v...", string(body), delay)
		time.Sleep(delay)
		SyncLaunchProvidersAndEvents(client)
		return
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Launches: HTTP %d: %s", resp.StatusCode, string(body))
		return
	}

	var result LaunchAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("❌ Invalid JSON: %v", err)
		return
	}

	for _, l := range result.Results {
		// Sync Provider
		provider := l.LaunchServiceProvider
		var providerPBID, rocketPBID, padPBID, missionPBID string
		if provider.ID != 0 && provider.Name != "" {
			providerRecord, _ := client.FindRecordByField("launch_providers", "spacedevs_id", provider.ID)
			if providerRecord == nil {
				created, err := client.CreateRecord("launch_providers", map[string]interface{}{
					"spacedevs_id":  provider.ID,
					"name":          provider.Name,
					"abbrev":        provider.Abbrev,
					"type":          provider.Type,
					"country_code":  provider.Country,
					"description":   provider.Desc,
					"founding_year": provider.Founded,
					"logo_url":      provider.LogoURL,
					"image_url":     provider.ImageURL,
					"wiki_url":      provider.WikiURL,
					"info_url":      provider.InfoURL,
				})
				if err != nil {
					log.Printf("❌ Failed to insert provider %s: %v", provider.Name, err)
				} else {
					providerPBID = (*created)["id"].(string)
				}
			} else {
				providerPBID = (*providerRecord)["id"].(string)
			}
		}

		// Sync Rocket
		if r := l.Rocket; r.ID != 0 && r.FullName != "" {
			rocketRecord, _ := client.FindRecordByField("rockets", "spacedevs_id", r.ID)
			if rocketRecord == nil {
				created, err := client.CreateRecord("rockets", map[string]interface{}{
					"spacedevs_id": r.ID,
					"name":         r.Name,
					"full_name":    r.FullName,
					"variant":      r.Variant,
					"family":       r.Family,
					"reusable":     r.Reusable,
					"description":  r.Description,
					"launch_mass":  r.LaunchMass,
					"leo_capacity": r.LEOCapacity,
					"gto_capacity": r.GTOCapacity,
					"image_url":    r.ImageURL,
					"info_url":     r.InfoURL,
					"wiki_url":     r.WikiURL,
					"manufacturer": r.Manufacturer.Name,
				})
				if err != nil {
					log.Printf("❌ Failed to insert rocket %s: %v", r.FullName, err)
				} else {
					rocketPBID = (*created)["id"].(string)
				}
			} else {
				rocketPBID = (*rocketRecord)["id"].(string)
			}
		}

		// Sync Pad
		if p := l.Pad; p.ID != 0 && p.Name != "" {
			padRecord, _ := client.FindRecordByField("pads", "spacedevs_id", p.ID)
			if padRecord == nil {
				created, err := client.CreateRecord("pads", map[string]interface{}{
					"spacedevs_id":       p.ID,
					"name":               p.Name,
					"description":        p.Description,
					"latitude":           p.Latitude,
					"longitude":          p.Longitude,
					"country_code":       p.CountryCode,
					"location_name":      p.LocationName,
					"map_url":            p.MapURL,
					"wiki_url":           p.WikiURL,
					"map_image":          p.MapImage,
					"launch_count_year":  p.LaunchCountYear,
					"launch_count_total": p.LaunchCountTotal,
				})
				if err != nil {
					log.Printf("❌ Failed to insert pad %s: %v", p.Name, err)
				} else {
					padPBID = (*created)["id"].(string)
				}
			} else {
				padPBID = (*padRecord)["id"].(string)
			}
		}

		// Sync Mission
		if m := l.Mission; m != nil && m.ID != 0 && m.Name != "" {
			missionRecord, _ := client.FindRecordByField("missions", "spacedevs_id", m.ID)
			if missionRecord == nil {
				created, err := client.CreateRecord("missions", map[string]interface{}{
					"spacedevs_id": m.ID,
					"name":         m.Name,
					"description":  m.Description,
					"type":         m.Type,
					"orbit":        m.Orbit.Name,
				})
				if err != nil {
					log.Printf("❌ Failed to insert mission %s: %v", m.Name, err)
				} else {
					missionPBID = (*created)["id"].(string)
				}
			} else {
				missionPBID = (*missionRecord)["id"].(string)
			}
		}

		// Sync Event
		if eventRecord, _ := client.FindRecordByField("events", "title", l.Name); eventRecord == nil {
			launchTime, _ := time.Parse(time.RFC3339, l.Net)
			// Convert []Update to []pbclient.Update
			var pbUpdates []pbclient.Update
			for _, u := range l.Updates {
				pbUpdates = append(pbUpdates, pbclient.Update{
					ID:          strconv.Itoa(u.ID),
					Title:       u.Comment,   // API 'comment' maps to our 'title'
					Description: u.InfoURL,   // API 'info_url' maps to our 'description'
					CreatedAt:   u.CreatedOn, // API 'created_on' maps to our 'created_at'
				})
			}
			err := client.CreateEvent(pbclient.Event{
				Title:       l.Name,
				Type:        "rocket_launch",
				Datetime:    launchTime.Format(time.RFC3339),
				Location:    l.Pad.Name,
				SourceURL:   l.URL,
				Description: "Synced from Launch Library",
				SpacedevsID: providerPBID, // must be PB ID for relation
				ProviderID:  providerPBID,
				RocketID:    rocketPBID,
				PadID:       padPBID,
				MissionID:   missionPBID,
				Updates:     pbUpdates,
			})
			if err != nil {
				log.Printf("❌ Failed to insert event %s: %v", l.Name, err)
			} else {
				log.Printf("✅ Synced event: %s", l.Name)
			}
		} else {
			log.Printf("⏭️ Event %s already exists", l.Name)
		}
	}

	log.Println("✅ All launches and related data synced.")
}
