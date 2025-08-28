//
//  MissionListView.swift
//  Station98
//
//  Created by Liam Arbuckle on 25/8/2025.
//

import SwiftUI

struct MissionsView: View {
    @State private var missions: [Mission] = []
    @State private var isLoading = true
    @State private var error: String?
    
    var body: some View {
        NavigationView {
            content
                .navigationTitle("🛰️ Missions")
        }
        .task {
            await loadMissions()
        }
    }
    
    @ViewBuilder
    private var content: some View {
        if isLoading {
            ProgressView("Loading missions...")
        } else if let error = error {
            Text("Error: \(error)")
                .foregroundColor(.red)
        } else if missions.isEmpty {
            Text("No missions available.")
                .foregroundColor(.gray)
        } else {
            List(missions) { mission in
                MissionCard(mission: mission)
            }
        }
    }
    
    private func loadMissions() async {
        let cacheFilename = "cached_missions.json"
        
        // Try loading from cache first
        if let cached: [Mission] = LocalCache.load(from: cacheFilename, as: [Mission].self), !cached.isEmpty {
            self.missions = cached
            self.isLoading = false
            return
        }

        do {
            let fetched = try await MissionFetcher.shared.fetchMissions()
            self.missions = fetched
            LocalCache.save(fetched, to: cacheFilename)
        } catch {
            self.error = error.localizedDescription
        }
        
        self.isLoading = false
    }
}
