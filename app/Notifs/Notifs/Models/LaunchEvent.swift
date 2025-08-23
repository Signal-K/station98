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

    var formattedDate: String {
        let isoFormatter = ISO8601DateFormatter()
        let dateFormatter = DateFormatter()
        dateFormatter.dateStyle = .medium
        dateFormatter.timeStyle = .short
        if let date = isoFormatter.date(from: datetime) {
            return dateFormatter.string(from: date)
        } else {
            return "Unknown date"
        }
    }
}
 
