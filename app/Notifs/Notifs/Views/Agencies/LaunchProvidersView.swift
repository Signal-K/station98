//
//  LaunchProvidersView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import Foundation
import SwiftUI

let cachedProviders: [LaunchProvider]? = LocalCache.load(from: "cached_providers.json", as: [LaunchProvider].self)

struct LaunchProvidersView: View {
    @StateObject var fetcher = LaunchProviderFetcher()

    var body: some View {
        NavigationView {
            List {
                if fetcher.isLoading {
                    ProgressView("Loading...")
                } else if let error = fetcher.error {
                    Text("Error: \(error)")
                        .foregroundColor(.red)
                } else {
                    let groupedProviders = Dictionary(grouping: fetcher.providers, by: { provider in
                        fetcher.allEvents.contains(where: { $0.spacedevs_id == String(provider.spacedevs_id ?? 0) })
                    })

                    let withEvents = groupedProviders[true] ?? []
                    let withoutEvents = groupedProviders[false] ?? []

                    if !withEvents.isEmpty {
                        Section(header: Text("🚀 Active Providers")) {
                            ForEach(withEvents) { provider in
                                NavigationLink(destination: ProviderEventsView(provider: provider, allEvents: fetcher.allEvents)) {
                                    Text(provider.name)
                                }
                            }
                        }
                    }

                    if !withoutEvents.isEmpty {
                        Section(header: Text("🛰️ Inactive Providers")) {
                            ForEach(withoutEvents) { provider in
                                NavigationLink(destination: ProviderEventsView(provider: provider, allEvents: fetcher.allEvents)) {
                                    Text(provider.name)
                                }
                            }
                        }
                    }
                }
            }
            .navigationTitle("Launch Providers")
            .task {
                if let cached = cachedProviders {
                    fetcher.providers = cached
                }
                await fetcher.fetch()
            }
            .overlay(content: {
                EmptyView()
            })
        }
    }
}

struct ProviderDetailView: View {
    let provider: LaunchProvider

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text(provider.name)
                .font(.largeTitle)
                .bold()

            if let description = provider.description {
                Text(description)
                    .font(.body)
            }

            if let country = provider.country_code {
                Text("Country: \(country)")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }

            Spacer()
        }
        .padding()
        .navigationTitle(provider.name)
    }
}
