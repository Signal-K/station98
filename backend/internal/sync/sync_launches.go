package sync

// SyncError is a simple error type for HTTP sync errors
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

// SyncError is a simple error type for HTTP sync errors
type SyncError struct {
	StatusCode int
	Body       string
}

func (e *SyncError) Error() string {
	return "Sync error: HTTP " + strconv.Itoa(e.StatusCode) + ": " + e.Body
}

func NewSyncError(code int, body string) error {
	return &SyncError{StatusCode: code, Body: body}
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
	go func() {
		offset := 0
		fetchCount := 0
		for {
			log.Printf("📡 Fetching launches (offset %d)...", offset)
			url := "https://ll.thespacedevs.com/2.2.0/launch/upcoming/?limit=50&mode=detailed&offset=" + strconv.Itoa(offset)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("❌ Failed to fetch launches: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			if resp.StatusCode == 429 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				delay := getThrottleDelay(string(body))
				log.Printf("❌ Launches: HTTP 429: %s. Retrying in %v...", string(body), delay)
				time.Sleep(delay)
				continue
			}
			if resp.StatusCode != 200 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				log.Printf("❌ Launches: HTTP %d: %s", resp.StatusCode, string(body))
				time.Sleep(30 * time.Minute)
				continue
			}

			var result LaunchAPIResponse
			err = json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()
			if err != nil {
				log.Printf("❌ Invalid JSON: %v", err)
				time.Sleep(30 * time.Minute)
				continue
			}

			if len(result.Results) == 0 {
				log.Printf("⏭️ No more launches found at offset %d. Resetting to offset 0.", offset)
				offset = 0
				fetchCount = 0
				time.Sleep(30 * time.Minute)
				continue
			}

			for _, l := range result.Results {
				// ...existing sync logic for provider, rocket, pad, mission, event...
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
							Title:       u.Comment,
							Description: u.InfoURL,
							CreatedAt:   u.CreatedOn,
						})
					}
					// Convert []Link (l.VideoURLs) to []map[string]interface{} for vid_urls
					var pbVidURLs []map[string]interface{}
					for _, v := range l.VideoURLs {
						pbVidURLs = append(pbVidURLs, map[string]interface{}{
							"title":    v.Title,
							"url":      v.URL,
							"priority": v.Priority,
						})
					}
					_, err := client.CreateRecord("events", map[string]interface{}{
						"title":        l.Name,
						"type":         "rocket_launch",
						"datetime":     launchTime.Format(time.RFC3339),
						"location":     l.Pad.Name,
						"source_url":   l.URL,
						"description":  "Synced from Launch Library",
						"spacedevs_id": providerPBID,
						"provider":     providerPBID,
						"rocket_id":    rocketPBID,
						"pad_id":       padPBID,
						"mission_id":   missionPBID,
						"updates":      pbUpdates,
						"vid_urls":     pbVidURLs,
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

			log.Println("✅ Launches and related data synced for offset", offset)
			offset += 50
			fetchCount++
			// Wait 1 minute after first fetch, then 30 minutes for subsequent fetches
			if fetchCount == 1 {
				time.Sleep(1 * time.Minute)
			} else {
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}
