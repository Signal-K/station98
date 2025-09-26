package migrations
package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}

		// Launch Window & Status Fields
		collection.Fields.Add(
			&core.DateTimeField{
				Name: "window_start",
			},
			&core.DateTimeField{
				Name: "window_end",
			},
			&core.URLField{
				Name: "infographic",
			},
			&core.BoolField{
				Name: "webcast_live",
			},
			&core.TextField{
				Name: "status_abbrev",
				Max:  10,
			},
			&core.TextField{
				Name: "status_description",
				Max:  500,
			},
		)

		// Enhanced Rocket Information
		collection.Fields.Add(
			&core.TextField{
				Name: "rocket_name",
				Max:  100,
			},
			&core.TextField{
				Name: "rocket_full_name",
				Max:  150,
			},
			&core.NumberField{
				Name: "rocket_total_launches",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "rocket_successful_launches",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "rocket_failed_launches",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "rocket_pending_launches",
				Min:  types.Pointer(0.0),
			},
		)

		// Vehicle/Booster Details
		collection.Fields.Add(
			&core.TextField{
				Name: "launcher_serial_number",
				Max:  50,
			},
			&core.NumberField{
				Name: "launcher_flight_number",
				Min:  types.Pointer(0.0),
			},
			&core.BoolField{
				Name: "launcher_reused",
			},
			&core.NumberField{
				Name: "launcher_flights",
				Min:  types.Pointer(0.0),
			},
			&core.TextField{
				Name: "launcher_status",
				Max:  50,
			},
		)

		// Landing Information
		collection.Fields.Add(
			&core.BoolField{
				Name: "landing_attempt",
			},
			&core.BoolField{
				Name: "landing_success",
			},
			&core.TextField{
				Name: "landing_location",
				Max:  100,
			},
			&core.TextField{
				Name: "landing_type",
				Max:  50,
			},
		)

		// Program Information
		collection.Fields.Add(
			&core.TextField{
				Name: "program_names",
				Max:  500,
			},
			&core.TextField{
				Name: "program_descriptions",
				Max:  2000,
			},
			&core.TextField{
				Name: "program_image_urls",
				Max:  1000,
			},
		)

		// Launch Statistics
		collection.Fields.Add(
			&core.NumberField{
				Name: "orbital_launch_attempt_count",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "location_launch_attempt_count",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "pad_launch_attempt_count",
				Min:  types.Pointer(0.0),
			},
			&core.NumberField{
				Name: "agency_launch_attempt_count",
				Min:  types.Pointer(0.0),
			},
		)

		// Crew Information
		collection.Fields.Add(
			&core.JSONField{
				Name: "crew_members",
			},
		)

		return app.Save(collection)
	}, func(app core.App) error {
		// Revert operation - remove the added fields
		collection, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return err
		}

		fieldsToRemove := []string{
			"window_start", "window_end", "infographic", "webcast_live",
			"status_abbrev", "status_description", "rocket_name", "rocket_full_name",
			"rocket_total_launches", "rocket_successful_launches", "rocket_failed_launches",
			"rocket_pending_launches", "launcher_serial_number", "launcher_flight_number",
			"launcher_reused", "launcher_flights", "launcher_status", "landing_attempt",
			"landing_success", "landing_location", "landing_type", "program_names",
			"program_descriptions", "program_image_urls", "orbital_launch_attempt_count",
			"location_launch_attempt_count", "pad_launch_attempt_count",
			"agency_launch_attempt_count", "crew_members",
		}

		for _, fieldName := range fieldsToRemove {
			collection.Fields.RemoveByName(fieldName)
		}

		return app.Save(collection)
	})
}