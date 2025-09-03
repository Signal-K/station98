//
//  Expedition.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import Foundation

struct Expedition: Identifiable, Codable {
    let id: String
    let name: String
    let startDate: String
    let endDate: String
    let patches: String?
    let expand: Expand?
    let station: String?

    var stationName: String {
        expand?.station.name ?? "Unknown"
    }

    struct Expand: Codable {
        let station: Station

        struct Station: Codable {
            let name: String
        }
    }

    var formattedDate: String {
        return "\(startDate) – \(endDate)"
    }

    enum CodingKeys: String, CodingKey {
        case id, name, patches, expand, station
        case startDate = "start_date"
        case endDate = "end_date"
    }
}
