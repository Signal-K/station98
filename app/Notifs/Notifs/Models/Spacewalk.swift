//
//  Spacewalk.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import Foundation

import Foundation

struct Spacewalk: Codable, Identifiable {
    let id: String
    let apiId: String
    let name: String
    let slug: String?
    let eventId: Int?
    let url: String?
    let location: String?
    let startTime: String
    let endTime: String?
    let duration: String?
    let expeditionName: String?

    enum CodingKeys: String, CodingKey {
        case id, name, slug, url, location, startTime = "start_time", endTime = "end_time", duration, expeditionName = "expedition_name", eventId = "event_id", apiId = "api_id"
    }

    var startDate: Date? {
        ISO8601DateFormatter().date(from: startTime)
    }

    var isArchived: Bool {
        guard let date = startDate else { return false }
        let calendar = Calendar.current
        let currentYear = calendar.component(.year, from: Date())
        let walkYear = calendar.component(.year, from: date)
        return walkYear < currentYear
    }
}
