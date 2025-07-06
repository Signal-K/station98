//
//  FetchSpacewalks.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import Foundation

@MainActor
class SpacewalksViewModel: ObservableObject {
    @Published var recent: [Spacewalk] = []
    @Published var archived: [Spacewalk] = []
    @Published var loading = true

    func fetchSpacewalks() async {
        guard let url = URL(string: "http://127.0.0.1:8080/api/collections/spacewalks/records?perPage=100&sort=-start_time") else { return }

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoded = try JSONDecoder().decode(SpacewalkListResponse.self, from: data)
            let all = decoded.items

            recent = all.filter { !$0.isArchived }
            archived = all.filter { $0.isArchived }
            loading = false
        } catch {
            print("❌ Error fetching spacewalks:", error)
            loading = false
        }
    }
}

struct SpacewalkListResponse: Codable {
    let items: [Spacewalk]
}
