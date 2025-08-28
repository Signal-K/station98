import Foundation

@MainActor
class EventFetcher: ObservableObject {
    @Published var events: [LaunchEvent] = []
    @Published var isLoading = false
    @Published var error: String?

    let baseURL = "http://localhost:8080"
    private let cacheFilename = "cached_events.json"

    func fetchEvents() async {
        isLoading = true
        error = nil
        
        // Try to fetch remote data first
        guard let url = URL(string: "\(baseURL)/api/collections/events/records?perPage=200") else {
            self.error = "Invalid URL"
            return
        }

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoded = try JSONDecoder().decode(EventResponse.self, from: data)
            self.events = decoded.items

            // Save to cache
            LocalCache.save(decoded, to: cacheFilename)
        } catch {
            // Try to load from cache if server fails
            if let cached: EventResponse = LocalCache.load(from: cacheFilename, as: EventResponse.self) {
                self.events = cached.items
                self.error = nil // clear error because we have fallback
            } else {
                self.error = "Could not connect to the server."
            }
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
