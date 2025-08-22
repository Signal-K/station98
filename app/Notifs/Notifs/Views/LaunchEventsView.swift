//
//  LaunchEventsView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import SwiftUI

struct LaunchEventsView: View {
    @StateObject private var fetcher = EventFetcher()
    @StateObject private var providerFetcher = LaunchProviderFetcher()

    var body: some View {
        NavigationView {
            content
                .navigationTitle("🚀 Launches")
        }
        .task {
            await fetcher.fetchEvents()
            await providerFetcher.fetch()
        }
    }

    @ViewBuilder
    private var content: some View {
        if fetcher.isLoading {
            ProgressView("Loading launches...")
        } else if let error = fetcher.error {
            Text("Error: \(error)")
                .foregroundColor(.red)
                .multilineTextAlignment(.center)
        } else {
            List(fetcher.events) { event in
                VStack(alignment: .leading, spacing: 4) {
                    Text(event.title)
                        .font(.headline)

                    Text(event.formattedDate)
                        .font(.subheadline)
                        .foregroundColor(.gray)

                    if let missionID = event.mission_id {
                        Text("Mission ID: \(missionID)")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }

                    if let providerIDString = event.spacedevs_id,
                       let providerID = Int(providerIDString),
                       let provider = providerFetcher.providers.first(where: { $0.spacedevs_id == providerID }) {
                        Text(provider.name)
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
            }
        }
    }
}
