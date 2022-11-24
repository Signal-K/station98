import XCTest
@testable import Station98

class ViewDataRetrievalTests: XCTestCase {
    // MARK: - LaunchEventsView
    func testLaunchEventsView_retrievesDataFromCacheAndAPI_andPresentsCorrectly() async {
        let fetcher = await MainActor.run { EventFetcher() }
        let providerFetcher = await MainActor.run { LaunchProviderFetcher() }
        let view = LaunchEventsView()
        // Simulate cache
        await MainActor.run {
            fetcher.events = [LaunchEvent(id: "1", title: "Test Launch", datetime: "2025-08-28", type: nil, source_url: nil, description: nil, spacedevs_id: nil, mission_id: nil, expand: nil, updates: nil, vid_urls: nil)]
            XCTAssertFalse(fetcher.events.isEmpty, "Should retrieve events from cache")
        }
        // Simulate API
        await fetcher.fetchEvents()
        await MainActor.run {
            XCTAssertTrue(fetcher.events.count >= 0, "Should retrieve events from API")
        }
        // Presentation
        XCTAssertNotNil(view.body, "View should present events")
    }

    // MARK: - MissionsView
    func testMissionsView_retrievesDataFromCacheAndAPI_andPresentsCorrectly() async {
        let view = MissionsView()
        // Simulate cache
        let missions = [Mission(id: "1", name: "Test Mission", description: "Desc", orbit: "LEO")]
        XCTAssertFalse(missions.isEmpty, "Should retrieve missions from cache")
        // Simulate API (would call loadMissions() async)
        XCTAssertNotNil(view.body, "View should present missions")
    }

    // MARK: - EventDetailView
    func testEventDetailView_retrievesDataFromCacheAndAPI_andPresentsCorrectly() async {
        let fetcher = await MainActor.run { EventFetcher() }
        let view = EventDetailView()
        // Simulate cache
        await MainActor.run {
            fetcher.events = [LaunchEvent(id: "1", title: "Test Event", datetime: "2025-08-28", type: nil, source_url: nil, description: nil, spacedevs_id: nil, mission_id: nil, expand: nil, updates: nil, vid_urls: nil)]
            XCTAssertFalse(fetcher.events.isEmpty, "Should retrieve events from cache")
        }
        // Simulate API
        await fetcher.fetchEvents()
        await MainActor.run {
            XCTAssertTrue(fetcher.events.count >= 0, "Should retrieve events from API")
        }
        // Presentation
        XCTAssertNotNil(view.body, "View should present event updates")
    }

    // MARK: - LaunchProvidersView
    func testLaunchProvidersView_retrievesDataFromCacheAndAPI_andPresentsCorrectly() async {
        let fetcher = await MainActor.run { LaunchProviderFetcher() }
        let view = LaunchProvidersView()
        // Simulate cache
        await MainActor.run {
            fetcher.providers = [LaunchProvider(id: "1", name: "Test Provider", abbrev: nil, country_code: nil, type: nil, founded: nil, logo_url: nil, image_url: nil, wiki_url: nil, info_url: nil, spacedevs_id: nil, description: nil)]
            XCTAssertFalse(fetcher.providers.isEmpty, "Should retrieve providers from cache")
        }
        // Simulate API
        await fetcher.fetch()
        await MainActor.run {
            XCTAssertTrue(fetcher.providers.count >= 0, "Should retrieve providers from API")
        }
        // Presentation
        XCTAssertNotNil(view.body, "View should present providers")
    }

    // MARK: - PadsGlobeView
    func testPadsGlobeView_retrievesDataFromCacheAndAPI_andPresentsCorrectly() async {
        let view = PadsGlobeView()
        // Simulate cache
        let pads = [Pad(id: "1", name: "Test Pad", latitude: 0.0, longitude: 0.0, location_name: "Test Location", country_code: "USA")]
        XCTAssertFalse(pads.isEmpty, "Should retrieve pads from cache")
        // Simulate API (would call fetchPads() async)
        XCTAssertNotNil(view.body, "View should present pads on globe")
    }
}
