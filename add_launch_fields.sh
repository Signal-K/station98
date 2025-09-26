#!/bin/bash

# Add enhanced launch tracking fields to the events table
echo "🔧 Adding enhanced launch tracking fields to events table..."

# SQL to add all the new columns
docker-compose exec pocketbase sqlite3 /pb/data/data.db << 'EOF'
-- Launch Window & Status Fields
ALTER TABLE events ADD COLUMN window_start TEXT;
ALTER TABLE events ADD COLUMN window_end TEXT;
ALTER TABLE events ADD COLUMN infographic TEXT;
ALTER TABLE events ADD COLUMN webcast_live INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN status_abbrev TEXT;
ALTER TABLE events ADD COLUMN status_description TEXT;

-- Enhanced Rocket Information
ALTER TABLE events ADD COLUMN rocket_name TEXT;
ALTER TABLE events ADD COLUMN rocket_full_name TEXT;
ALTER TABLE events ADD COLUMN rocket_total_launches INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN rocket_successful_launches INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN rocket_failed_launches INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN rocket_pending_launches INTEGER DEFAULT 0;

-- Vehicle/Booster Details
ALTER TABLE events ADD COLUMN launcher_serial_number TEXT;
ALTER TABLE events ADD COLUMN launcher_flight_number INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN launcher_reused INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN launcher_flights INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN launcher_status TEXT;

-- Landing Information
ALTER TABLE events ADD COLUMN landing_attempt INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN landing_success INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN landing_location TEXT;
ALTER TABLE events ADD COLUMN landing_type TEXT;

-- Program Information
ALTER TABLE events ADD COLUMN program_names TEXT;
ALTER TABLE events ADD COLUMN program_descriptions TEXT;
ALTER TABLE events ADD COLUMN program_image_urls TEXT;

-- Launch Statistics
ALTER TABLE events ADD COLUMN orbital_launch_attempt_count INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN location_launch_attempt_count INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN pad_launch_attempt_count INTEGER DEFAULT 0;
ALTER TABLE events ADD COLUMN agency_launch_attempt_count INTEGER DEFAULT 0;

-- Crew Information
ALTER TABLE events ADD COLUMN crew_members TEXT; -- JSON as TEXT in SQLite
EOF

echo "✅ Enhanced launch tracking fields added!"
echo "🔄 Restart the backend to start populating these fields with launch data."