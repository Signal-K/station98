/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("events")

  // Launch Window & Status Fields
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "window_start",
    "name": "window_start",
    "type": "date",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "window_end", 
    "name": "window_end",
    "type": "date",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "infographic",
    "name": "infographic", 
    "type": "url",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "exceptDomains": null,
      "onlyDomains": null
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "webcast_live",
    "name": "webcast_live",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "status_abbrev",
    "name": "status_abbrev",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 10,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "status_description",
    "name": "status_description",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 500,
      "pattern": ""
    }
  }))

  // Enhanced Rocket Information
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_name",
    "name": "rocket_name",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 100,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_full_name",
    "name": "rocket_full_name",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 150,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_total_launches",
    "name": "rocket_total_launches",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_successful_launches",
    "name": "rocket_successful_launches",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_failed_launches",
    "name": "rocket_failed_launches",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rocket_pending_launches",
    "name": "rocket_pending_launches",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  // Vehicle/Booster Details
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "launcher_serial_number",
    "name": "launcher_serial_number",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 50,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "launcher_flight_number",
    "name": "launcher_flight_number",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "launcher_reused",
    "name": "launcher_reused",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "launcher_flights",
    "name": "launcher_flights",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "launcher_status",
    "name": "launcher_status",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 50,
      "pattern": ""
    }
  }))

  // Landing Information
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "landing_attempt",
    "name": "landing_attempt",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "landing_success",
    "name": "landing_success",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "landing_location",
    "name": "landing_location",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 100,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "landing_type",
    "name": "landing_type",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 50,
      "pattern": ""
    }
  }))

  // Program Information
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "program_names",
    "name": "program_names",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 500,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "program_descriptions",
    "name": "program_descriptions",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 2000,
      "pattern": ""
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "program_image_urls",
    "name": "program_image_urls",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": 1000,
      "pattern": ""
    }
  }))

  // Launch Statistics
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "orbital_launch_attempt_count",
    "name": "orbital_launch_attempt_count",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "location_launch_attempt_count",
    "name": "location_launch_attempt_count",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "pad_launch_attempt_count",
    "name": "pad_launch_attempt_count",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "agency_launch_attempt_count",
    "name": "agency_launch_attempt_count",
    "type": "number",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": 0,
      "max": null,
      "noDecimal": true
    }
  }))

  // Crew Information
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "crew_members",
    "name": "crew_members",
    "type": "json",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSize": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("events")

  // Remove all the added fields (rollback)
  collection.schema.removeField("window_start")
  collection.schema.removeField("window_end")
  collection.schema.removeField("infographic")
  collection.schema.removeField("webcast_live")
  collection.schema.removeField("status_abbrev")
  collection.schema.removeField("status_description")
  collection.schema.removeField("rocket_name")
  collection.schema.removeField("rocket_full_name")
  collection.schema.removeField("rocket_total_launches")
  collection.schema.removeField("rocket_successful_launches")
  collection.schema.removeField("rocket_failed_launches")
  collection.schema.removeField("rocket_pending_launches")
  collection.schema.removeField("launcher_serial_number")
  collection.schema.removeField("launcher_flight_number")
  collection.schema.removeField("launcher_reused")
  collection.schema.removeField("launcher_flights")
  collection.schema.removeField("launcher_status")
  collection.schema.removeField("landing_attempt")
  collection.schema.removeField("landing_success")
  collection.schema.removeField("landing_location")
  collection.schema.removeField("landing_type")
  collection.schema.removeField("program_names")
  collection.schema.removeField("program_descriptions")
  collection.schema.removeField("program_image_urls")
  collection.schema.removeField("orbital_launch_attempt_count")
  collection.schema.removeField("location_launch_attempt_count")
  collection.schema.removeField("pad_launch_attempt_count")
  collection.schema.removeField("agency_launch_attempt_count")
  collection.schema.removeField("crew_members")

  return dao.saveCollection(collection)
})