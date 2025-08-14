//
//  ContentView.swift
//  SpaceNotifs
//
//  Created by Liam Arbuckle on 15/8/2025.
//

import SwiftUI

struct ContentView: View {
    @StateObject private var fetcher = EventFetcher()

    var body: some View {
        NavigationView {
            Group {
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
                            if let location = event.location, !location.isEmpty {
                                Text("📍 \(location)")
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                        }
                    }
                }
            }
            .navigationTitle("🚀 Launches")
        }
        .task {
            await fetcher.fetchEvents()
        }
    }
}
