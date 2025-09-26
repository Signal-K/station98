#!/usr/bin/env node

const axios = require('axios');

const POCKETBASE_URL = 'http://localhost:8080';

async function addFieldsToEventsCollection() {
    try {
        console.log('🔧 Adding enhanced launch tracking fields to events collection...');
        
        // First, let's get the current events collection schema
        console.log('📋 Fetching current events collection schema...');
        
        // Get collections without auth first to see what we're working with
        try {
            const collectionsResponse = await axios.get(`${POCKETBASE_URL}/api/collections`);
            console.log('✅ Collections accessible, found', collectionsResponse.data.items.length, 'collections');
            
            // Look for events collection
            const eventsCollection = collectionsResponse.data.items.find(c => c.name === 'events');
            if (eventsCollection) {
                console.log('✅ Found events collection with', eventsCollection.schema.length, 'existing fields');
                console.log('📝 Current fields:', eventsCollection.schema.map(f => f.name).join(', '));
            } else {
                console.log('❌ Events collection not found');
                return;
            }
        } catch (error) {
            console.log('❌ Cannot access collections:', error.response?.status, error.response?.statusText);
            console.log('💡 This is expected if authentication is required');
        }
        
        console.log(`
📋 MANUAL STEPS REQUIRED:
Since PocketBase admin authentication is not working via API, please follow these steps:

1. Open http://localhost:8080/_/ in your browser
2. Login to PocketBase admin (or create admin account if first time)
3. Go to Collections > events
4. Add the following fields using "New field" button:

LAUNCH WINDOW & STATUS:
- window_start (DateTime) 
- window_end (DateTime)
- infographic (URL)
- webcast_live (Bool)
- status_abbrev (Text, max 10)
- status_description (Text, max 500)

ROCKET INFORMATION:
- rocket_name (Text, max 100)
- rocket_full_name (Text, max 150) 
- rocket_total_launches (Number, min 0)
- rocket_successful_launches (Number, min 0)
- rocket_failed_launches (Number, min 0)
- rocket_pending_launches (Number, min 0)

VEHICLE/BOOSTER:
- launcher_serial_number (Text, max 50)
- launcher_flight_number (Number, min 0)
- launcher_reused (Bool)
- launcher_flights (Number, min 0)
- launcher_status (Text, max 50)

LANDING INFO:
- landing_attempt (Bool)
- landing_success (Bool) 
- landing_location (Text, max 100)
- landing_type (Text, max 50)

PROGRAM INFO:
- program_names (Text, max 500)
- program_descriptions (Text, max 2000)
- program_image_urls (Text, max 1000)

STATISTICS:
- orbital_launch_attempt_count (Number, min 0)
- location_launch_attempt_count (Number, min 0) 
- pad_launch_attempt_count (Number, min 0)
- agency_launch_attempt_count (Number, min 0)

CREW:
- crew_members (JSON)

5. Save the collection
6. Restart the backend service to begin syncing enhanced data

⚠️  IMPORTANT: Only ADD fields - do not remove existing ones to preserve data!
        `);

    } catch (error) {
        console.error('❌ Error:', error.message);
        if (error.response) {
            console.error('Response:', error.response.status, error.response.statusText);
            console.error('Data:', error.response.data);
        }
    }
}

if (require.main === module) {
    addFieldsToEventsCollection();
}

module.exports = { addFieldsToEventsCollection };