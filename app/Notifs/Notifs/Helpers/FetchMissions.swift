//
//  FetchMissions.swift
//  Station98
//
//  Created by Liam Arbuckle on 25/8/2025.
//

import Foundation

struct PocketbaseResponse<T: Codable>: Codable {
    let items: [T]
}

class MissionFetcher {
    static let shared = MissionFetcher()
    private let baseURL = "http://localhost:8080"
    private let cacheFilename = "cached_missions.json"

    func fetchMissions() async throws -> [Mission] {
        if let cached: PocketbaseResponse<Mission> = LocalCache.load(from: cacheFilename, as: PocketbaseResponse<Mission>.self) {
            return cached.items
        }

        guard let url = URL(string: "\(baseURL)/api/collections/missions/records?perPage=200") else {
            throw URLError(.badURL)
        }

        let (data, _) = try await URLSession.shared.data(from: url)
        let decoded = try JSONDecoder().decode(PocketbaseResponse<Mission>.self, from: data)
        LocalCache.save(decoded, to: cacheFilename)
        return decoded.items
    }
}
