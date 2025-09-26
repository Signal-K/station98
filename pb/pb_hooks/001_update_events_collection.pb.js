// PocketBase JavaScript Migration Script
// File: pb_hooks/001_update_events_collection.pb.js
// This script will run automatically when PocketBase starts

migrate((txApp) => {
    const collection = txApp.findCollectionByNameOrId("events");
    if (!collection) {
        throw new Error("Events collection not found");
    }

    // Add Launch Window & Status Fields
    collection.fields.add(new DateField({
        name: "window_start",
        presentable: false,
    }));

    collection.fields.add(new DateField({
        name: "window_end", 
        presentable: false,
    }));

    collection.fields.add(new URLField({
        name: "infographic",
        presentable: false,
    }));

    collection.fields.add(new BoolField({
        name: "webcast_live",
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "status_abbrev",
        max: 10,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "status_description",
        max: 500,
        presentable: false,
    }));

    // Add Enhanced Rocket Information
    collection.fields.add(new TextField({
        name: "rocket_name",
        max: 100,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "rocket_full_name",
        max: 150,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "rocket_total_launches",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "rocket_successful_launches",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "rocket_failed_launches",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "rocket_pending_launches",
        min: 0,
        presentable: false,
    }));

    // Add Vehicle/Booster Details
    collection.fields.add(new TextField({
        name: "launcher_serial_number",
        max: 50,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "launcher_flight_number",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new BoolField({
        name: "launcher_reused",
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "launcher_flights",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "launcher_status",
        max: 50,
        presentable: false,
    }));

    // Add Landing Information
    collection.fields.add(new BoolField({
        name: "landing_attempt",
        presentable: false,
    }));

    collection.fields.add(new BoolField({
        name: "landing_success",
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "landing_location",
        max: 100,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "landing_type",
        max: 50,
        presentable: false,
    }));

    // Add Program Information
    collection.fields.add(new TextField({
        name: "program_names",
        max: 500,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "program_descriptions",
        max: 2000,
        presentable: false,
    }));

    collection.fields.add(new TextField({
        name: "program_image_urls",
        max: 1000,
        presentable: false,
    }));

    // Add Launch Statistics
    collection.fields.add(new NumberField({
        name: "orbital_launch_attempt_count",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "location_launch_attempt_count",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "pad_launch_attempt_count",
        min: 0,
        presentable: false,
    }));

    collection.fields.add(new NumberField({
        name: "agency_launch_attempt_count",
        min: 0,
        presentable: false,
    }));

    // Add Crew Information (JSON field for crew member details)
    collection.fields.add(new JSONField({
        name: "crew_members",
        presentable: false,
    }));

    return txApp.save(collection);
}, (txApp) => {
    // Downgrade - remove the fields we added
    const collection = txApp.findCollectionByNameOrId("events");
    if (!collection) {
        return; // Collection doesn't exist, nothing to downgrade
    }

    // List of field names to remove during downgrade
    const fieldsToRemove = [
        "window_start", "window_end", "infographic", "webcast_live", 
        "status_abbrev", "status_description", "rocket_name", "rocket_full_name",
        "rocket_total_launches", "rocket_successful_launches", "rocket_failed_launches",
        "rocket_pending_launches", "launcher_serial_number", "launcher_flight_number",
        "launcher_reused", "launcher_flights", "launcher_status", "landing_attempt",
        "landing_success", "landing_location", "landing_type", "program_names",
        "program_descriptions", "program_image_urls", "orbital_launch_attempt_count",
        "location_launch_attempt_count", "pad_launch_attempt_count", 
        "agency_launch_attempt_count", "crew_members"
    ];

    fieldsToRemove.forEach(fieldName => {
        collection.fields.removeById(fieldName);
    });

    return txApp.save(collection);
});