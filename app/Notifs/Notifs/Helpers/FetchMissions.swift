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

    func fetchMissions() async throws -> [Mission] {
        guard let url = URL(string: "\(baseURL)/api/collections/missions/records?perPage=200") else {
            throw URLError(.badURL)
        }

        let (data, _) = try await URLSession.shared.data(from: url)
        let decoded = try JSONDecoder().decode(PocketbaseResponse<Mission>.self, from: data)
        return decoded.items
    }
}
