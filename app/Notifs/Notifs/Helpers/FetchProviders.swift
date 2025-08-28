//
//  FetchProviders.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import Foundation

@MainActor
class LaunchProviderFetcher: ObservableObject {
    @Published var providers: [LaunchProvider] = []
    @Published var filteredProviders: [LaunchProvider] = []
    @Published var allEvents: [LaunchEvent] = []
    @Published var isLoading = false
    @Published var error: String?

    func fetch() async {
        isLoading = true
        error = nil

        var cachedProviders: LaunchProviderResult?
        var cachedEvents: LaunchEventResult?

        if let loadedProviders: LaunchProviderResult = LocalCache.load(from: "cached_providers.json", as: LaunchProviderResult.self) {
            cachedProviders = loadedProviders
            self.providers = loadedProviders.items
        }

        if let loadedEvents: LaunchEventResult = LocalCache.load(from: "cached_events_for_providers.json", as: LaunchEventResult.self) {
            cachedEvents = loadedEvents
            self.allEvents = loadedEvents.items
        }

        if let cachedProviders = cachedProviders, let cachedEvents = cachedEvents {
            let activeIDs = Set(cachedEvents.items.map { $0.spacedevs_id })
            self.filteredProviders = cachedProviders.items.filter { provider in
                if let id = provider.spacedevs_id {
                    return activeIDs.contains(String(id))
                }
                return false
            }
        }

        do {
            // Fetch providers
            let providerURL = URL(string: "http://localhost:8080/api/collections/launch_providers/records?perPage=200")!
            let (providerData, _) = try await URLSession.shared.data(from: providerURL)
            let providerResult = try JSONDecoder().decode(LaunchProviderResult.self, from: providerData)
            self.providers = providerResult.items
            LocalCache.save(providerResult, to: "cached_providers.json")

            // Fetch events
            let eventURL = URL(string: "http://localhost:8080/api/collections/events/records?perPage=200")!
            let (eventData, _) = try await URLSession.shared.data(from: eventURL)
            let eventResult = try JSONDecoder().decode(LaunchEventResult.self, from: eventData)
            self.allEvents = eventResult.items
            LocalCache.save(eventResult, to: "cached_events_for_providers.json")

            // Filter providers that have current launches
            let activeIDs = Set(eventResult.items.map { $0.spacedevs_id })
            self.filteredProviders = providerResult.items.filter { provider in
                if let id = provider.spacedevs_id {
                    return activeIDs.contains(String(id))
                }
                return false
            }
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }
}

struct LaunchProviderResult: Codable {
    let items: [LaunchProvider]
}

struct LaunchEventResult: Codable {
    let items: [LaunchEvent]
}
