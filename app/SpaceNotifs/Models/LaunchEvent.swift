import Foundation

struct EventResponse: Codable {
    let items: [LaunchEvent]
}

struct LaunchEvent: Identifiable, Codable {
    let id: String
    let title: String
    let datetime: String
    let location: String?
    let sourceURL: String?

    var formattedDate: String {
        let formatter = ISO8601DateFormatter()
        if let date = formatter.date(from: datetime) {
            let displayFormatter = DateFormatter()
            displayFormatter.dateStyle = .medium
            displayFormatter.timeStyle = .short
            return displayFormatter.string(from: date)
        }
        return datetime
    }
} 