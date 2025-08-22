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

class MissionEventFetcher {
    static let shared = MissionEventFetcher()

    func fetchMissionsAndLinkedEvents(completion: @escaping ([(Mission, [LaunchEvent])]) -> Void) {
        guard let missionURL = URL(string: "http://localhost:8080/api/collections/missions/records?perPage=200"),
              let eventURL = URL(string: "http://localhost:8080/api/collections/events/records?perPage=200&expand=mission_id") else {
            completion([])
            return
        }

        var missions: [Mission] = []
        var events: [LaunchEvent] = []

        let group = DispatchGroup()

        // Fetch Missions
        group.enter()
        URLSession.shared.dataTask(with: missionURL) { data, _, _ in
            defer { group.leave() }
            guard let data = data else { return }
            do {
                let decoded = try JSONDecoder().decode(PocketbaseResponse<Mission>.self, from: data)
                missions = decoded.items
            } catch {
                print("Failed to decode missions: \(error)")
            }
        }.resume()

        // Fetch Events
        group.enter()
        URLSession.shared.dataTask(with: eventURL) { data, _, _ in
            defer { group.leave() }
            guard let data = data else { return }
            do {
                let decoded = try JSONDecoder().decode(PocketbaseResponse<LaunchEvent>.self, from: data)
                events = decoded.items.filter { $0.expand?.mission != nil }
            } catch {
                print("Failed to decode events: \(error)")
            }
        }.resume()

        // When both are done
        group.notify(queue: .main) {
            let groupedEvents = Dictionary(grouping: events, by: { $0.expand!.mission!.id })

            let result = missions.compactMap { mission -> (Mission, [LaunchEvent])? in
                guard let matchingEvents = groupedEvents[mission.id] else {
                    return nil
                }
                return (mission, matchingEvents)
            }

            completion(result)
        }
    }
}
