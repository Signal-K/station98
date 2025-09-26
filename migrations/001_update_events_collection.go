package main

// Migration script to add enhanced launch tracking fields to events collection
// This file contains the field definitions for manual application via PocketBase admin UI

/*
MIGRATION: Enhanced Launch Tracking Fields for Events Collection

This migration adds comprehensive launch information fields to the events collection
without removing or deleting existing rows.

FIELDS TO ADD:

1. Launch Window & Status Fields:
   - window_start (DateTime) - Launch window start time
   - window_end (DateTime) - Launch window end time
   - infographic (URL) - Link to launch infographic
   - webcast_live (Bool) - Whether webcast is live
   - status_abbrev (Text, max 10) - Launch status abbreviation
   - status_description (Text, max 500) - Full status description

2. Enhanced Rocket Information:
   - rocket_name (Text, max 100) - Rocket configuration name
   - rocket_full_name (Text, max 150) - Full rocket name
   - rocket_total_launches (Number, min 0) - Total launches for this rocket type
   - rocket_successful_launches (Number, min 0) - Successful launches count
   - rocket_failed_launches (Number, min 0) - Failed launches count
   - rocket_pending_launches (Number, min 0) - Pending launches count

3. Vehicle/Booster Details:
   - launcher_serial_number (Text, max 50) - Booster serial number
   - launcher_flight_number (Number, min 0) - Flight number for this booster
   - launcher_reused (Bool) - Whether booster is reused
   - launcher_flights (Number, min 0) - Total flights for this booster
   - launcher_status (Text, max 50) - Booster status

4. Landing Information:
   - landing_attempt (Bool) - Whether landing will be attempted
   - landing_success (Bool) - Whether landing was successful
   - landing_location (Text, max 100) - Landing location name
   - landing_type (Text, max 50) - Landing type (e.g., "ASDS", "RTLS")

5. Program Information:
   - program_names (Text, max 500) - Space program names (comma-separated)
   - program_descriptions (Text, max 2000) - Program descriptions (pipe-separated)
   - program_image_urls (Text, max 1000) - Program image URLs (comma-separated)

6. Launch Statistics:
   - orbital_launch_attempt_count (Number, min 0) - Overall orbital launch attempt count
   - location_launch_attempt_count (Number, min 0) - Launch attempts from this location
   - pad_launch_attempt_count (Number, min 0) - Launch attempts from this pad
   - agency_launch_attempt_count (Number, min 0) - Launch attempts by this agency

7. Crew Information:
   - crew_members (JSON) - Array of crew member details for crewed missions

INSTRUCTIONS:
1. Open PocketBase admin dashboard (http://localhost:8080/_/)
2. Go to Collections > events
3. Add each field listed above using the "New field" button
4. Set the field type and constraints as specified
5. Save the collection

IMPORTANT: This migration is designed to ADD fields only. Existing data will be preserved.
*/
