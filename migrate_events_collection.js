#!/usr/bin/env node

/**
 * PocketBase Collection Migration Script
 * 
 * This script adds enhanced launch tracking fields to the events collection
 * via the PocketBase Admin API without affecting existing data.
 * 
 * Usage:
 * 1. Make sure PocketBase is running (docker-compose up pocketbase)
 * 2. Install dependencies: npm install axios
 * 3. Run: node migrate_events_collection.js
 */

const axios = require('axios');

const POCKETBASE_URL = 'http://localhost:8080';
const ADMIN_EMAIL = 'teddy@scroobl.es';
const ADMIN_PASSWORD = 'teddy@scroobl.es';

async function authenticateAdmin() {
    try {
        const response = await axios.post(`${POCKETBASE_URL}/api/admins/auth-with-password`, {
            identity: ADMIN_EMAIL,
            password: ADMIN_PASSWORD
        });
        
        return response.data.token;
    } catch (error) {
        console.error('Failed to authenticate admin:', error.response?.data || error.message);
        throw error;
    }
}

async function getEventsCollection(token) {
    try {
        const response = await axios.get(`${POCKETBASE_URL}/api/collections/events`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('Failed to get events collection:', error.response?.data || error.message);
        throw error;
    }
}

async function updateEventsCollection(token, collection) {
    // Define the new fields to add
    const newFields = [
        // Launch Window & Status Fields
        { name: "window_start", type: "date", presentable: false },
        { name: "window_end", type: "date", presentable: false },
        { name: "infographic", type: "url", presentable: false },
        { name: "webcast_live", type: "bool", presentable: false },
        { name: "status_abbrev", type: "text", max: 10, presentable: false },
        { name: "status_description", type: "text", max: 500, presentable: false },
        
        // Enhanced Rocket Information
        { name: "rocket_name", type: "text", max: 100, presentable: false },
        { name: "rocket_full_name", type: "text", max: 150, presentable: false },
        { name: "rocket_total_launches", type: "number", min: 0, presentable: false },
        { name: "rocket_successful_launches", type: "number", min: 0, presentable: false },
        { name: "rocket_failed_launches", type: "number", min: 0, presentable: false },
        { name: "rocket_pending_launches", type: "number", min: 0, presentable: false },
        
        // Vehicle/Booster Details
        { name: "launcher_serial_number", type: "text", max: 50, presentable: false },
        { name: "launcher_flight_number", type: "number", min: 0, presentable: false },
        { name: "launcher_reused", type: "bool", presentable: false },
        { name: "launcher_flights", type: "number", min: 0, presentable: false },
        { name: "launcher_status", type: "text", max: 50, presentable: false },
        
        // Landing Information
        { name: "landing_attempt", type: "bool", presentable: false },
        { name: "landing_success", type: "bool", presentable: false },
        { name: "landing_location", type: "text", max: 100, presentable: false },
        { name: "landing_type", type: "text", max: 50, presentable: false },
        
        // Program Information
        { name: "program_names", type: "text", max: 500, presentable: false },
        { name: "program_descriptions", type: "text", max: 2000, presentable: false },
        { name: "program_image_urls", type: "text", max: 1000, presentable: false },
        
        // Launch Statistics
        { name: "orbital_launch_attempt_count", type: "number", min: 0, presentable: false },
        { name: "location_launch_attempt_count", type: "number", min: 0, presentable: false },
        { name: "pad_launch_attempt_count", type: "number", min: 0, presentable: false },
        { name: "agency_launch_attempt_count", type: "number", min: 0, presentable: false },
        
        // Crew Information
        { name: "crew_members", type: "json", presentable: false }
    ];

    // Check which fields already exist to avoid duplicates
    const existingFieldNames = collection.fields.map(field => field.name);
    const fieldsToAdd = newFields.filter(field => !existingFieldNames.includes(field.name));

    if (fieldsToAdd.length === 0) {
        console.log('All fields already exist. No migration needed.');
        return;
    }

    console.log(`Adding ${fieldsToAdd.length} new fields to events collection...`);

    // Add new fields to the collection
    collection.fields = [...collection.fields, ...fieldsToAdd];

    try {
        const response = await axios.patch(`${POCKETBASE_URL}/api/collections/${collection.id}`, collection, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });
        
        console.log('✅ Successfully updated events collection with enhanced launch tracking fields');
        console.log(`Added fields: ${fieldsToAdd.map(f => f.name).join(', ')}`);
        
        return response.data;
    } catch (error) {
        console.error('Failed to update events collection:', error.response?.data || error.message);
        throw error;
    }
}

async function main() {
    try {
        console.log('🚀 Starting PocketBase events collection migration...');
        
        // Authenticate as admin
        console.log('🔐 Authenticating admin...');
        const token = await authenticateAdmin();
        
        // Get current events collection
        console.log('📋 Fetching events collection...');
        const collection = await getEventsCollection(token);
        
        // Update the collection with new fields
        console.log('🔄 Updating collection with enhanced fields...');
        await updateEventsCollection(token, collection);
        
        console.log('✅ Migration completed successfully!');
        console.log('');
        console.log('The events collection now supports comprehensive launch tracking including:');
        console.log('- Launch window times and status information');
        console.log('- Detailed rocket configuration and performance metrics');
        console.log('- Vehicle/booster reuse and flight history');
        console.log('- Landing attempt and success tracking');
        console.log('- Space program associations');
        console.log('- Launch statistics and crew information');
        console.log('');
        console.log('Existing event records are preserved and will be populated with the new data on next sync.');
        
    } catch (error) {
        console.error('❌ Migration failed:', error.message);
        process.exit(1);
    }
}

// Run the migration
main();