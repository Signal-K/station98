//
//  NotifsApp.swift
//  Notifs
//
//  Created by Liam Arbuckle on 16/8/2025.
//

import SwiftUI
import SwiftData
import Clerk

@main
struct NotifsApp: App {
    @State private var clerk = Clerk.shared
    
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
                .environment(\.clerk, clerk)
                .task {
                  clerk.configure(publishableKey: "pk_test_d2VsY29tZWQtZG9nZmlzaC0zLmNsZXJrLmFjY291bnRzLmRldiQ")
                  try? await clerk.load()
                }
        }
        .modelContainer(sharedModelContainer)
    }
}
