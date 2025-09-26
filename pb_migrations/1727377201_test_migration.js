/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("events")

  // Add just a few test fields first to verify the migration works
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "test_window_start",
    "name": "window_start",
    "type": "date",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "test_webcast_live",
    "name": "webcast_live",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "test_rocket_name",
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

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("events")

  // Remove test fields
  collection.schema.removeField("test_window_start")
  collection.schema.removeField("test_webcast_live")  
  collection.schema.removeField("test_rocket_name")

  return dao.saveCollection(collection)
})