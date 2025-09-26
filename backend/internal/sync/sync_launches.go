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
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	Slug                       string          `json:"slug"`
	Net                        string          `json:"net"`
	WindowStart                string          `json:"window_start"`
	WindowEnd                  string          `json:"window_end"`
	URL                        string          `json:"url"`
	Status                     Status          `json:"status"`
	LaunchServiceProvider      Agency          `json:"launch_service_provider"`
	Mission                    *Mission        `json:"mission"`
	Rocket                     Rocket          `json:"rocket"`
	Pad                        Pad             `json:"pad"`
	InfoURLs                   []Link          `json:"infoURLs"`
	VideoURLs                  []Link          `json:"vidURLs"`
	WebcastLive                bool            `json:"webcast_live"`
	Image                      string          `json:"image"`
	Infographic                string          `json:"infographic"`
	Timeline                   []TimelineEntry `json:"timeline"`
	Updates                    []Update        `json:"updates"`
	Program                    []Program       `json:"program"`
	OrbitalLaunchAttemptCount  int             `json:"orbital_launch_attempt_count"`
	LocationLaunchAttemptCount int             `json:"location_launch_attempt_count"`
	PadLaunchAttemptCount      int             `json:"pad_launch_attempt_count"`
	AgencyLaunchAttemptCount   int             `json:"agency_launch_attempt_count"`
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

type Program struct {
	ID             int            `json:"id"`
	URL            string         `json:"url"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Agencies       []Agency       `json:"agencies"`
	ImageURL       string         `json:"image_url"`
	StartDate      string         `json:"start_date"`
	EndDate        *string        `json:"end_date"`
	InfoURL        *string        `json:"info_url"`
	WikiURL        string         `json:"wiki_url"`
	MissionPatches []MissionPatch `json:"mission_patches"`
	Type           ProgramType    `json:"type"`
}

type MissionPatch struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	ImageURL string `json:"image_url"`
	Agency   Agency `json:"agency"`
}

type ProgramType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Rocket struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	FullName        string           `json:"full_name"`
	Variant         string           `json:"variant"`
	Family          string           `json:"family"`
	Reusable        bool             `json:"reusable"`
	Description     string           `json:"description"`
	LaunchMass      int              `json:"launch_mass"`
	LEOCapacity     int              `json:"leo_capacity"`
	GTOCapacity     int              `json:"gto_capacity"`
	ImageURL        string           `json:"image_url"`
	InfoURL         string           `json:"info_url"`
	WikiURL         string           `json:"wiki_url"`
	Manufacturer    Agency           `json:"manufacturer"`
	Configuration   RocketConfig     `json:"configuration"`
	LauncherStage   []LauncherStage  `json:"launcher_stage"`
	SpacecraftStage *SpacecraftStage `json:"spacecraft_stage"`
}

type RocketConfig struct {
	ID                            int       `json:"id"`
	Name                          string    `json:"name"`
	FullName                      string    `json:"full_name"`
	Family                        string    `json:"family"`
	Variant                       string    `json:"variant"`
	Active                        bool      `json:"active"`
	Reusable                      bool      `json:"reusable"`
	Description                   string    `json:"description"`
	Manufacturer                  Agency    `json:"manufacturer"`
	Program                       []Program `json:"program"`
	ImageURL                      string    `json:"image_url"`
	InfoURL                       string    `json:"info_url"`
	WikiURL                       string    `json:"wiki_url"`
	TotalLaunchCount              int       `json:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int       `json:"consecutive_successful_launches"`
	SuccessfulLaunches            int       `json:"successful_launches"`
	FailedLaunches                int       `json:"failed_launches"`
	PendingLaunches               int       `json:"pending_launches"`
}

type LauncherStage struct {
	ID                   int         `json:"id"`
	Type                 string      `json:"type"`
	Reused               *bool       `json:"reused"`
	LauncherFlightNumber *int        `json:"launcher_flight_number"`
	Launcher             Launcher    `json:"launcher"`
	Landing              *Landing    `json:"landing"`
	PreviousFlightDate   *string     `json:"previous_flight_date"`
	TurnAroundTimeDays   *int        `json:"turn_around_time_days"`
	PreviousFlight       interface{} `json:"previous_flight"`
}

type Launcher struct {
	ID                 int     `json:"id"`
	URL                string  `json:"url"`
	Details            string  `json:"details"`
	FlightProven       bool    `json:"flight_proven"`
	SerialNumber       string  `json:"serial_number"`
	Status             string  `json:"status"`
	ImageURL           string  `json:"image_url"`
	SuccessfulLandings *int    `json:"successful_landings"`
	AttemptedLandings  *int    `json:"attempted_landings"`
	Flights            *int    `json:"flights"`
	LastLaunchDate     *string `json:"last_launch_date"`
	FirstLaunchDate    *string `json:"first_launch_date"`
}

type Landing struct {
	ID                int             `json:"id"`
	Attempt           bool            `json:"attempt"`
	Success           *bool           `json:"success"`
	Description       string          `json:"description"`
	DownrangeDistance *float64        `json:"downrange_distance"`
	Location          LandingLocation `json:"location"`
	Type              LandingType     `json:"type"`
}

type LandingLocation struct {
	ID                 int         `json:"id"`
	Name               string      `json:"name"`
	Abbrev             string      `json:"abbrev"`
	Description        string      `json:"description"`
	Location           interface{} `json:"location"`
	SuccessfulLandings int         `json:"successful_landings"`
}

type LandingType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type SpacecraftStage struct {
	ID              int          `json:"id"`
	URL             string       `json:"url"`
	MissionEnd      *string      `json:"mission_end"`
	Destination     string       `json:"destination"`
	LaunchCrew      []CrewMember `json:"launch_crew"`
	OnsiteCrew      []CrewMember `json:"onsite_crew"`
	LandingCrew     []CrewMember `json:"landing_crew"`
	Spacecraft      Spacecraft   `json:"spacecraft"`
	LaunchingTo     *string      `json:"launching_to"`
	LandingLocation *string      `json:"landing_location"`
}

type CrewMember struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Astronaut Astronaut `json:"astronaut"`
}

type Role struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Priority int    `json:"priority"`
}

type Astronaut struct {
	ID                    int             `json:"id"`
	Name                  string          `json:"name"`
	Status                AstronautStatus `json:"status"`
	Type                  AstronautType   `json:"type"`
	Agency                Agency          `json:"agency"`
	DateOfBirth           string          `json:"date_of_birth"`
	DateOfDeath           *string         `json:"date_of_death"`
	Nationality           string          `json:"nationality"`
	Bio                   string          `json:"bio"`
	TwitterURL            string          `json:"twitter"`
	InstagramURL          string          `json:"instagram"`
	WikiURL               string          `json:"wiki"`
	ProfileImageURL       string          `json:"profile_image"`
	ProfileImageThumbnail string          `json:"profile_image_thumbnail"`
}

type AstronautStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AstronautType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Spacecraft struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	SerialNumber     string           `json:"serial_number"`
	Status           SpacecraftStatus `json:"status"`
	Description      string           `json:"description"`
	SpacecraftConfig SpacecraftConfig `json:"spacecraft_config"`
}

type SpacecraftStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SpacecraftConfig struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Type            SpacecraftType `json:"type"`
	Agency          Agency         `json:"agency"`
	InUse           bool           `json:"in_use"`
	Capability      string         `json:"capability"`
	History         string         `json:"history"`
	Details         string         `json:"details"`
	MaidenFlight    string         `json:"maiden_flight"`
	Height          *float64       `json:"height"`
	Diameter        *float64       `json:"diameter"`
	HumanRated      bool           `json:"human_rated"`
	CrewCapacity    *int           `json:"crew_capacity"`
	PayloadCapacity *int           `json:"payload_capacity"`
	FlightLife      *string        `json:"flight_life"`
	ImageURL        string         `json:"image_url"`
	NationURL       string         `json:"nation_url"`
	WikiURL         string         `json:"wiki_url"`
	InfoURL         string         `json:"info_url"`
}

type SpacecraftType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
					windowStart, _ := time.Parse(time.RFC3339, l.WindowStart)
					windowEnd, _ := time.Parse(time.RFC3339, l.WindowEnd)

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

					// Convert []Link (l.InfoURLs) to []map[string]interface{} for info_urls
					var pbInfoURLs []map[string]interface{}
					for _, info := range l.InfoURLs {
						pbInfoURLs = append(pbInfoURLs, map[string]interface{}{
							"title":    info.Title,
							"url":      info.URL,
							"priority": info.Priority,
						})
					}

					// Convert Timeline to []map[string]interface{}
					var pbTimeline []map[string]interface{}
					for _, t := range l.Timeline {
						pbTimeline = append(pbTimeline, map[string]interface{}{
							"time":  t.Time,
							"event": t.Event,
						})
					}

					// Extract rocket configuration details
					rocketConfigName := ""
					rocketConfigFullName := ""
					rocketTotalLaunches := 0
					rocketSuccessfulLaunches := 0
					rocketFailedLaunches := 0
					rocketPendingLaunches := 0
					if l.Rocket.Configuration.ID != 0 {
						rocketConfigName = l.Rocket.Configuration.Name
						rocketConfigFullName = l.Rocket.Configuration.FullName
						rocketTotalLaunches = l.Rocket.Configuration.TotalLaunchCount
						rocketSuccessfulLaunches = l.Rocket.Configuration.SuccessfulLaunches
						rocketFailedLaunches = l.Rocket.Configuration.FailedLaunches
						rocketPendingLaunches = l.Rocket.Configuration.PendingLaunches
					}

					// Extract launcher stage details (first stage info)
					launcherSerialNumber := ""
					launcherFlightNumber := 0
					launcherReused := false
					launcherFlights := 0
					launcherStatus := ""
					landingAttempt := false
					landingSuccess := false
					landingLocation := ""
					landingType := ""
					if len(l.Rocket.LauncherStage) > 0 {
						stage := l.Rocket.LauncherStage[0]
						launcherSerialNumber = stage.Launcher.SerialNumber
						launcherStatus = stage.Launcher.Status
						if stage.LauncherFlightNumber != nil {
							launcherFlightNumber = *stage.LauncherFlightNumber
						}
						if stage.Reused != nil {
							launcherReused = *stage.Reused
						}
						if stage.Launcher.Flights != nil {
							launcherFlights = *stage.Launcher.Flights
						}
						if stage.Landing != nil {
							landingAttempt = stage.Landing.Attempt
							if stage.Landing.Success != nil {
								landingSuccess = *stage.Landing.Success
							}
							landingLocation = stage.Landing.Location.Name
							landingType = stage.Landing.Type.Name
						}
					}

					// Extract program information
					var programNames []string
					var programDescriptions []string
					var programImageURLs []string
					for _, prog := range l.Program {
						programNames = append(programNames, prog.Name)
						programDescriptions = append(programDescriptions, prog.Description)
						programImageURLs = append(programImageURLs, prog.ImageURL)
					}

					// Extract crew information if spacecraft stage exists
					var crewMembers []map[string]interface{}
					if l.Rocket.SpacecraftStage != nil {
						for _, crew := range l.Rocket.SpacecraftStage.LaunchCrew {
							crewMembers = append(crewMembers, map[string]interface{}{
								"astronaut_id":  crew.Astronaut.ID,
								"name":          crew.Astronaut.Name,
								"role":          crew.Role.Role,
								"role_priority": crew.Role.Priority,
								"nationality":   crew.Astronaut.Nationality,
								"agency":        crew.Astronaut.Agency.Name,
								"profile_image": crew.Astronaut.ProfileImageURL,
							})
						}
					}

					_, err := client.CreateRecord("events", map[string]interface{}{
						"title":                         l.Name,
						"type":                          "rocket_launch",
						"datetime":                      launchTime.Format(time.RFC3339),
						"window_start":                  windowStart.Format(time.RFC3339),
						"window_end":                    windowEnd.Format(time.RFC3339),
						"location":                      l.Pad.Name,
						"source_url":                    l.URL,
						"description":                   "Synced from Launch Library",
						"spacedevs_id":                  providerPBID,
						"provider":                      providerPBID,
						"rocket_id":                     rocketPBID,
						"pad_id":                        padPBID,
						"mission_id":                    missionPBID,
						"updates":                       pbUpdates,
						"vid_urls":                      pbVidURLs,
						"info_urls":                     pbInfoURLs,
						"timeline":                      pbTimeline,
						"image":                         l.Image,
						"infographic":                   l.Infographic,
						"webcast_live":                  l.WebcastLive,
						"status_abbrev":                 l.Status.Abbrev,
						"status_description":            l.Status.Description,
						"rocket_name":                   rocketConfigName,
						"rocket_full_name":              rocketConfigFullName,
						"rocket_total_launches":         rocketTotalLaunches,
						"rocket_successful_launches":    rocketSuccessfulLaunches,
						"rocket_failed_launches":        rocketFailedLaunches,
						"rocket_pending_launches":       rocketPendingLaunches,
						"launcher_serial_number":        launcherSerialNumber,
						"launcher_flight_number":        launcherFlightNumber,
						"launcher_reused":               launcherReused,
						"launcher_flights":              launcherFlights,
						"launcher_status":               launcherStatus,
						"landing_attempt":               landingAttempt,
						"landing_success":               landingSuccess,
						"landing_location":              landingLocation,
						"landing_type":                  landingType,
						"program_names":                 strings.Join(programNames, ", "),
						"program_descriptions":          strings.Join(programDescriptions, " | "),
						"program_image_urls":            strings.Join(programImageURLs, ", "),
						"crew_members":                  crewMembers,
						"orbital_launch_attempt_count":  l.OrbitalLaunchAttemptCount,
						"location_launch_attempt_count": l.LocationLaunchAttemptCount,
						"pad_launch_attempt_count":      l.PadLaunchAttemptCount,
						"agency_launch_attempt_count":   l.AgencyLaunchAttemptCount,
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
