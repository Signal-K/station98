import Foundation

@MainActor
class EventFetcher: ObservableObject {
    @Published var events: [LaunchEvent] = []
    @Published var isLoading = false
    @Published var error: String?

    let baseURL = "http://localhost:8080"

    func fetchEvents() async {
        isLoading = true
        error = nil

        guard let url = URL(string: "\(baseURL)/api/collections/events/records?filter=(mission_id!='')") else {
            self.error = "Invalid URL"
            return
        }

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoded = try JSONDecoder().decode(EventResponse.self, from: data)
            self.events = decoded.items
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }
    
    func fetchEventsWithVideosOrUpdates() async {
        isLoading = true
        error = nil

        guard let url = URL(string: "\(baseURL)/api/collections/events/records?filter=(vid_urls!=''||updates!='')") else {
            self.error = "Invalid URL"
            return
        }

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoded = try JSONDecoder().decode(EventResponse.self, from: data)
            self.events = decoded.items
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    struct EventResponse: Codable {
        let items: [LaunchEvent]
    }
}
