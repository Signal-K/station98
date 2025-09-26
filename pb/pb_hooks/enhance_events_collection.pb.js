/// <reference path="../data/types.d.ts" />
onServe((e) => {
  // Add enhanced launch fields to events collection on startup
  const dao = $app.dao()
  
  try {
    const collection = dao.findCollectionByNameOrId("events")
    let modified = false
    
    // Define all the fields we want to add
    const fieldsToAdd = [
      // Launch Window & Status
      { name: "window_start", type: "date" },
      { name: "window_end", type: "date" },
      { name: "infographic", type: "url" },
      { name: "webcast_live", type: "bool" },
      { name: "status_abbrev", type: "text", max: 10 },
      { name: "status_description", type: "text", max: 500 },
      
      // Rocket Information
      { name: "rocket_name", type: "text", max: 100 },
      { name: "rocket_full_name", type: "text", max: 150 },
      { name: "rocket_total_launches", type: "number", min: 0 },
      { name: "rocket_successful_launches", type: "number", min: 0 },
      { name: "rocket_failed_launches", type: "number", min: 0 },
      { name: "rocket_pending_launches", type: "number", min: 0 },
      
      // Vehicle/Booster Details
      { name: "launcher_serial_number", type: "text", max: 50 },
      { name: "launcher_flight_number", type: "number", min: 0 },
      { name: "launcher_reused", type: "bool" },
      { name: "launcher_flights", type: "number", min: 0 },
      { name: "launcher_status", type: "text", max: 50 },
      
      // Landing Information
      { name: "landing_attempt", type: "bool" },
      { name: "landing_success", type: "bool" },
      { name: "landing_location", type: "text", max: 100 },
      { name: "landing_type", type: "text", max: 50 },
      
      // Program Information
      { name: "program_names", type: "text", max: 500 },
      { name: "program_descriptions", type: "text", max: 2000 },
      { name: "program_image_urls", type: "text", max: 1000 },
      
      // Launch Statistics
      { name: "orbital_launch_attempt_count", type: "number", min: 0 },
      { name: "location_launch_attempt_count", type: "number", min: 0 },
      { name: "pad_launch_attempt_count", type: "number", min: 0 },
      { name: "agency_launch_attempt_count", type: "number", min: 0 },
      
      // Crew Information
      { name: "crew_members", type: "json" }
    ]
    
    // Check and add missing fields
    fieldsToAdd.forEach(fieldDef => {
      const existingField = collection.schema.find(f => f.name === fieldDef.name)
      if (!existingField) {
        console.log(`Adding field: ${fieldDef.name}`)
        
        const field = new SchemaField({
          system: false,
          id: fieldDef.name + "_" + Math.random().toString(36).substr(2, 9),
          name: fieldDef.name,
          type: fieldDef.type,
          required: false,
          presentable: false,
          unique: false,
          options: {}
        })
        
        // Set type-specific options
        if (fieldDef.type === "text" && fieldDef.max) {
          field.options.max = fieldDef.max
        } else if (fieldDef.type === "number" && fieldDef.min !== undefined) {
          field.options.min = fieldDef.min
          field.options.noDecimal = true
        }
        
        collection.schema.addField(field)
        modified = true
      }
    })
    
    if (modified) {
      dao.saveCollection(collection)
      console.log("✅ Enhanced launch tracking fields added to events collection")
    } else {
      console.log("✅ All enhanced launch tracking fields already exist")
    }
    
  } catch (error) {
    console.log("❌ Error adding enhanced launch fields:", error)
  }
})