//
//  NotifsApp.swift
//  Notifs
//
//  Created by Liam Arbuckle on 16/8/2025.
//

import SwiftUI
import SwiftData

import Appwrite

let client = Client()
    .setEndpoint("http://localhost:8020/v1")
    .setProject("station126")
    .setSelfSigned(true) 

@main
struct NotifsApp: App {
    var sharedModelContainer: ModelContainer = {
        let schema = Schema([
            Item.self,
        ])
        let modelConfiguration = ModelConfiguration(schema: schema, isStoredInMemoryOnly: false)

        do {
            return try ModelContainer(for: schema, configurations: [modelConfiguration])
        } catch {
            fatalError("Could not create ModelContainer: \(error)")
        } 
    }()

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
        .modelContainer(sharedModelContainer)
    }
}
