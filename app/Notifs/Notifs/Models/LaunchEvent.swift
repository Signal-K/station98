import Foundation

struct EventResponse: Codable {
    let items: [LaunchEvent]
}

struct LaunchEvent: Codable, Identifiable {
    let id: String
    let title: String
    let datetime: String
    let type: String?
    let source_url: String?
    let description: String?
    let spacedevs_id: String?
    let mission_id: String?
    let expand: LaunchEventExpand?
    
    struct LaunchEventExpand: Codable {
        let mission: Mission?
    }

    var formattedDate: String {
        let isoFormatter = ISO8601DateFormatter()
        isoFormatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]

        if let date = isoFormatter.date(from: datetime) {
            let formatter = DateFormatter()
            formatter.dateStyle = .medium
            formatter.timeStyle = .short
            return formatter.string(from: date)
        } else {
            return "Unknown date"
        }
    }
}
