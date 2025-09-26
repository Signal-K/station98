## Space Notifications: Enhanced Launch Tracking Setup

### Current Status
✅ **Backend Code Enhanced**: The `sync_launches.go` file has been updated with comprehensive launch tracking features including:
- Rocket configuration details (name, launch counts, success rates)
- Vehicle information (booster serial numbers, flight numbers, reuse status)
- Landing attempt/success tracking with location and type
- Space program associations 
- Crew member information for crewed flights
- Launch statistics and attempt counts
- Enhanced timeline and updates from SpaceDevs API

✅ **JSON Parsing Issue Fixed**: Fixed struct definition for `LauncherStage.PreviousFlight` to handle API response format

⏳ **Database Schema Update Needed**: The PocketBase `events` collection needs new fields added to store the enhanced launch data.

### Next Steps Required

#### 1. Add Fields to PocketBase Events Collection

**Manual Steps** (since admin API auth isn't working):

1. **Open PocketBase Admin**: Navigate to http://localhost:8080/_/
2. **Login/Create Admin Account** (if first time)
3. **Go to Collections**: Click on "events" collection
4. **Add New Fields**: Use the "New field" button to add each field below:

#### Required Fields to Add:

**Launch Window & Status:**
- `window_start` (DateTime) - Launch window start time
- `window_end` (DateTime) - Launch window end time  
- `infographic` (URL) - Link to launch infographic
- `webcast_live` (Bool) - Whether webcast is live
- `status_abbrev` (Text, max 10) - Launch status abbreviation
- `status_description` (Text, max 500) - Full status description

**Enhanced Rocket Information:**
- `rocket_name` (Text, max 100) - Rocket configuration name
- `rocket_full_name` (Text, max 150) - Full rocket name
- `rocket_total_launches` (Number, min 0) - Total launches for this rocket type
- `rocket_successful_launches` (Number, min 0) - Successful launches count
- `rocket_failed_launches` (Number, min 0) - Failed launches count
- `rocket_pending_launches` (Number, min 0) - Pending launches count

**Vehicle/Booster Details:**
- `launcher_serial_number` (Text, max 50) - Booster serial number
- `launcher_flight_number` (Number, min 0) - Flight number for this booster
- `launcher_reused` (Bool) - Whether booster is reused
- `launcher_flights` (Number, min 0) - Total flights for this booster
- `launcher_status` (Text, max 50) - Booster status

**Landing Information:**
- `landing_attempt` (Bool) - Whether landing will be attempted
- `landing_success` (Bool) - Whether landing was successful
- `landing_location` (Text, max 100) - Landing location name
- `landing_type` (Text, max 50) - Landing type (e.g., "ASDS", "RTLS")

**Program Information:**
- `program_names` (Text, max 500) - Space program names (comma-separated)
- `program_descriptions` (Text, max 2000) - Program descriptions (pipe-separated)
- `program_image_urls` (Text, max 1000) - Program image URLs (comma-separated)

**Launch Statistics:**
- `orbital_launch_attempt_count` (Number, min 0) - Overall orbital launch attempt count
- `location_launch_attempt_count` (Number, min 0) - Launch attempts from this location
- `pad_launch_attempt_count` (Number, min 0) - Launch attempts from this pad
- `agency_launch_attempt_count` (Number, min 0) - Launch attempts by this agency

**Crew Information:**
- `crew_members` (JSON) - Array of crew member details for crewed missions

#### 2. Restart Backend Service

Once fields are added to the collection:
```bash
docker-compose up backend
```

#### 3. Verify Enhanced Data

After restart, check that new events contain the enhanced fields:
- Rocket names and launch numbers
- Vehicle reuse information
- Landing attempt/success data
- Program associations
- Crew information (for crewed flights)
- Launch statistics

### What This Enables

**For Users:**
- See which booster is flying and if it's been used before
- Track landing attempts and success rates
- Understand space program context (Artemis, ISS, etc.)
- View crew information for human spaceflight missions
- Access comprehensive launch statistics

**For Developers:**
- Rich data for building advanced filtering and sorting
- Support for reuse tracking and statistics
- Program-based event categorization
- Crew tracking for human spaceflight features

### Important Notes

⚠️ **Data Preservation**: This migration only ADDS fields - existing event data is preserved
⚠️ **Order Matters**: Add fields to database BEFORE restarting backend to avoid sync errors  
⚠️ **Comprehensive Coverage**: New events will include 25+ additional data points from SpaceDevs API

### Files Modified

- `backend/internal/sync/sync_launches.go` - Enhanced with comprehensive launch tracking
- `migrations/001_update_events_collection.go` - Field specifications for manual application
- Various migration scripts created for different deployment scenarios

The backend is now ready to capture and store comprehensive launch tracking data once the database schema is updated!