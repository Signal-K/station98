package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type SpacewalkAPIResponse struct {
	Results []Spacewalk `json:"results"`
}

type Role struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Priority int    `json:"priority"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AgencyType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Agency struct {
	ResponseMode string     `json:"response_mode"`
	ID           int        `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	Abbrev       string     `json:"abbrev"`
	Type         AgencyType `json:"type"`
}

type Image struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Credit       string `json:"credit"`
	SingleUse    bool   `json:"single_use"`
}

type Nationality struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Alpha2Code             string `json:"alpha_2_code"`
	Alpha3Code             string `json:"alpha_3_code"`
	NationalityName        string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type AstronautType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Astronaut struct {
	ID               int             `json:"id"`
	URL              string          `json:"url"`
	Name             string          `json:"name"`
	Status           Status          `json:"status"`
	Agency           Agency          `json:"agency"`
	Image            *Image          `json:"image"`
	ResponseMode     string          `json:"response_mode"`
	Type             AstronautType   `json:"type"`
	InSpace          bool            `json:"in_space"`
	TimeInSpace      string          `json:"time_in_space"`
	EVATime          string          `json:"eva_time"`
	Age              int             `json:"age"`
	DateOfBirth      string          `json:"date_of_birth"`
	DateOfDeath      *string         `json:"date_of_death"`
	Nationality      []Nationality   `json:"nationality"`
	Bio              string          `json:"bio"`
	Wiki             *string         `json:"wiki"`
	LastFlight       string          `json:"last_flight"`
	FirstFlight      string          `json:"first_flight"`
}

type CrewMember struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Astronaut Astronaut `json:"astronaut"`
}

type SpaceStationType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SpaceStation struct {
	ID          int              `json:"id"`
	URL         string           `json:"url"`
	Name        string           `json:"name"`
	Image       *Image           `json:"image"`
	Status      Status           `json:"status"`
	Founded     string           `json:"founded"`
	Deorbited   *string          `json:"deorbited"`
	Description string           `json:"description"`
	Orbit       string           `json:"orbit"`
	Type        SpaceStationType `json:"type"`
}

type MissionPatch struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	ImageURL string `json:"image_url"`
	Agency   Agency `json:"agency"`
}

type Expedition struct {
	ID             int             `json:"id"`
	URL            string          `json:"url"`
	Name           string          `json:"name"`
	Start          string          `json:"start"`
	End            *string         `json:"end"`
	SpaceStation   SpaceStation    `json:"spacestation"`
	MissionPatches []MissionPatch  `json:"mission_patches"`
}

type ProgramType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Program struct {
	ResponseMode    string        `json:"response_mode"`
	ID              int           `json:"id"`
	URL             string        `json:"url"`
	Name            string        `json:"name"`
	Image           *Image        `json:"image"`
	InfoURL         *string       `json:"info_url"`
	WikiURL         *string       `json:"wiki_url"`
	Description     string        `json:"description"`
	Agencies        []Agency      `json:"agencies"`
	StartDate       string        `json:"start_date"`
	EndDate         *string       `json:"end_date"`
	MissionPatches  []interface{} `json:"mission_patches"`
	Type            ProgramType   `json:"type"`
}

type Spacewalk struct {
	ResponseMode     string      `json:"response_mode"`
	ID               int         `json:"id"`
	URL              string      `json:"url"`
	Name             string      `json:"name"`
	Start            string      `json:"start"`
	End              string      `json:"end"`
	Duration         string      `json:"duration"`
	Location         string      `json:"location"`
	Crew             []CrewMember `json:"crew"`
	SpaceStation     SpaceStation `json:"spacestation"`
	Expedition       Expedition   `json:"expedition"`
	SpacecraftFlight interface{}  `json:"spacecraft_flight"`
	Event            interface{}  `json:"event"`
	Program          []Program    `json:"program"`
}

func main() {
	fmt.Println("🚀 Fetching most recent spacewalk from Space Devs API...")

	// Hit the Space Devs API for spacewalks (order by ID descending to get most recent)
	url := "https://ll.thespacedevs.com/2.3.0/spacewalks/?limit=1&ordering=-id&mode=detailed&format=json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch spacewalks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result SpacewalkAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(result.Results) == 0 {
		fmt.Println("❌ No spacewalks found")
		return
	}

	// Display the first (most recent) spacewalk
	spacewalk := result.Results[0]
	
	fmt.Printf("✅ Most recent spacewalk found!\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════════════\n\n")
	
	// Basic Spacewalk Info
	fmt.Printf("🚶‍♂️ SPACEWALK EVENT\n")
	fmt.Printf("├─ ID: %d\n", spacewalk.ID)
	fmt.Printf("├─ Name: %s\n", spacewalk.Name)
	fmt.Printf("├─ URL: %s\n", spacewalk.URL)
	fmt.Printf("├─ Location: %s\n", spacewalk.Location)
	fmt.Printf("├─ Start Time: %s\n", formatDate(spacewalk.Start))
	fmt.Printf("├─ End Time: %s\n", formatDate(spacewalk.End))
	fmt.Printf("└─ Duration: %s\n", spacewalk.Duration)
	
	// Crew Information (Highlighted)
	fmt.Printf("\n👨‍🚀 SPACEWALK CREW\n")
	for i, crew := range spacewalk.Crew {
		prefix := "├─"
		if i == len(spacewalk.Crew)-1 {
			prefix = "└─"
		}
		fmt.Printf("%s Crew Member %d:\n", prefix, i+1)
		fmt.Printf("   ├─ Name: %s\n", crew.Astronaut.Name)
		fmt.Printf("   ├─ Role: %s (Priority: %d)\n", crew.Role.Role, crew.Role.Priority)
		fmt.Printf("   ├─ Age: %d\n", crew.Astronaut.Age)
		fmt.Printf("   ├─ Status: %s\n", crew.Astronaut.Status.Name)
		fmt.Printf("   ├─ In Space: %t\n", crew.Astronaut.InSpace)
		fmt.Printf("   ├─ Time in Space: %s\n", crew.Astronaut.TimeInSpace)
		fmt.Printf("   ├─ Total EVA Time: %s\n", crew.Astronaut.EVATime)
		fmt.Printf("   ├─ Date of Birth: %s\n", formatDate(crew.Astronaut.DateOfBirth))
		
		// Nationality
		if len(crew.Astronaut.Nationality) > 0 {
			nationalities := make([]string, len(crew.Astronaut.Nationality))
			for j, nat := range crew.Astronaut.Nationality {
				nationalities[j] = fmt.Sprintf("%s (%s)", nat.NationalityName, nat.Alpha2Code)
			}
			fmt.Printf("   ├─ Nationality: %s\n", strings.Join(nationalities, ", "))
		}
		
		fmt.Printf("   ├─ First Flight: %s\n", formatDate(crew.Astronaut.FirstFlight))
		fmt.Printf("   ├─ Last Flight: %s\n", formatDate(crew.Astronaut.LastFlight))
		
		// Agency
		fmt.Printf("   ├─ Agency: %s (%s)\n", crew.Astronaut.Agency.Name, crew.Astronaut.Agency.Abbrev)
		fmt.Printf("   ├─ Agency Type: %s\n", crew.Astronaut.Agency.Type.Name)
		
		if crew.Astronaut.Bio != "" {
			bioPreview := crew.Astronaut.Bio
			if len(bioPreview) > 200 {
				bioPreview = bioPreview[:200] + "..."
			}
			fmt.Printf("   ├─ Bio: %s\n", bioPreview)
		}
		
		if crew.Astronaut.Wiki != nil && *crew.Astronaut.Wiki != "" {
			fmt.Printf("   └─ Wikipedia: %s\n", *crew.Astronaut.Wiki)
		} else {
			fmt.Printf("   └─ Wikipedia: Not available\n")
		}
	}
	
	// Space Station Information (Highlighted)
	fmt.Printf("\n🛰 SPACE STATION\n")
	fmt.Printf("├─ Name: %s\n", spacewalk.SpaceStation.Name)
	fmt.Printf("├─ Status: %s\n", spacewalk.SpaceStation.Status.Name)
	fmt.Printf("├─ Type: %s\n", spacewalk.SpaceStation.Type.Name)
	fmt.Printf("├─ Founded: %s\n", formatDate(spacewalk.SpaceStation.Founded))
	fmt.Printf("├─ Orbit: %s\n", spacewalk.SpaceStation.Orbit)
	fmt.Printf("├─ Description: %s\n", spacewalk.SpaceStation.Description)
	fmt.Printf("└─ URL: %s\n", spacewalk.SpaceStation.URL)
	
	// Expedition Information (Highlighted)
	fmt.Printf("\n🚀 EXPEDITION\n")
	fmt.Printf("├─ Name: %s\n", spacewalk.Expedition.Name)
	fmt.Printf("├─ Start Date: %s\n", formatDate(spacewalk.Expedition.Start))
	if spacewalk.Expedition.End != nil {
		fmt.Printf("├─ End Date: %s\n", formatDate(*spacewalk.Expedition.End))
	} else {
		fmt.Printf("├─ End Date: Ongoing\n")
	}
	fmt.Printf("├─ Space Station: %s\n", spacewalk.Expedition.SpaceStation.Name)
	fmt.Printf("├─ URL: %s\n", spacewalk.Expedition.URL)
	
	// Mission Patches
	if len(spacewalk.Expedition.MissionPatches) > 0 {
		fmt.Printf("└─ Mission Patches:\n")
		for i, patch := range spacewalk.Expedition.MissionPatches {
			patchPrefix := "   ├─"
			if i == len(spacewalk.Expedition.MissionPatches)-1 {
				patchPrefix = "   └─"
			}
			fmt.Printf("%s %s (Priority: %d) - %s\n", patchPrefix, patch.Name, patch.Priority, patch.ImageURL)
		}
	} else {
		fmt.Printf("└─ Mission Patches: None\n")
	}
	
	// Program Information (Highlighted)
	if len(spacewalk.Program) > 0 {
		fmt.Printf("\n🌌 PROGRAMS\n")
		for i, program := range spacewalk.Program {
			prefix := "├─"
			if i == len(spacewalk.Program)-1 {
				prefix = "└─"
			}
			fmt.Printf("%s Program %d:\n", prefix, i+1)
			fmt.Printf("   ├─ Name: %s\n", program.Name)
			fmt.Printf("   ├─ Type: %s\n", program.Type.Name)
			fmt.Printf("   ├─ Description: %s\n", program.Description)
			fmt.Printf("   ├─ Start Date: %s\n", formatDate(program.StartDate))
			if program.EndDate != nil {
				fmt.Printf("   ├─ End Date: %s\n", formatDate(*program.EndDate))
			} else {
				fmt.Printf("   ├─ End Date: Ongoing\n")
			}
			
			if len(program.Agencies) > 0 {
				fmt.Printf("   ├─ Agencies:\n")
				for j, agency := range program.Agencies {
					agencyPrefix := "   │  ├─"
					if j == len(program.Agencies)-1 {
						agencyPrefix = "   │  └─"
					}
					fmt.Printf("%s %s (%s) - %s\n", agencyPrefix, agency.Name, agency.Abbrev, agency.Type.Name)
				}
			}
			
			if program.WikiURL != nil && *program.WikiURL != "" {
				fmt.Printf("   ├─ Wikipedia: %s\n", *program.WikiURL)
			}
			
			fmt.Printf("   └─ URL: %s\n", program.URL)
		}
	}
	
	fmt.Printf("\n═══════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("🎯 Spacewalk data retrieved successfully!\n")
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return "Not available"
	}
	return dateStr
}
