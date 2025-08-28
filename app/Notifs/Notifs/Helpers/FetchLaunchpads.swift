//
//  FetchLaunchpads.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import Foundation

@MainActor
class PadFetcher: ObservableObject {
    @Published var pads: [Pad] = []
    @Published var isLoading = false
    @Published var error: String?

    func fetchPads() async {
        isLoading = true
        error = nil

        // Try loading from cache first
        if let cached: PadResponse = LocalCache.load(from: "cached_pads.json", as: PadResponse.self) {
            self.pads = cached.items
        }

        guard let url = URL(string: "http://localhost:8080/api/collections/pads/records?perPage=200") else {
            error = "Invalid URL"
            return
        }

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoded = try JSONDecoder().decode(PadResponse.self, from: data)
            pads = decoded.items
            
            // Save to cache
            LocalCache.save(decoded, to: "cached_pads.json")
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }
}
