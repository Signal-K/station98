//
//  LaunchEventsView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import SwiftUI

struct LaunchEventsView: View {
    @StateObject private var fetcher = EventFetcher()
    @StateObject private var providerFetcher = LaunchProviderFetcher()
    
    var body: some View {
        NavigationView {
            content
                .navigationTitle("🎸 Launches")
        }
        .task {
            await fetcher.fetchEvents()
            await providerFetcher.fetch()
        }
    }
    
    @ViewBuilder
    private var content: some View {
        if fetcher.isLoading {
            ProgressView("Loading launches...")
        } else if let error = fetcher.error {
            Text("Error: \(error)")
                .foregroundColor(.red)
                .multilineTextAlignment(.center)
        } else {
            let grouped = groupEvents(events: fetcher.events)
            
            List {
                if !grouped.today.isEmpty {
                    Section(header: Text("🌞 Today")) {
                        ForEach(grouped.today) { event in
                            EventRow(event: event)
                        }
                    }
                }
                
                if !grouped.thisWeek.isEmpty {
                    Section(header: Text("📆 This Week")) {
                        ForEach(grouped.thisWeek) { event in
                            EventRow(event: event)
                        }
                    }
                }
                
                if !grouped.future.isEmpty {
                    Section(header: Text("🔭 Future")) {
                        ForEach(grouped.future) { event in
                            EventRow(event: event)
                        }
                    }
                }
                
                if !grouped.past.isEmpty {
                    Section(header: Text("📜 Past")) {
                        ForEach(grouped.past) { event in
                            EventRow(event: event)
                        }
                    }
                }
            }
        }
    }
    
    private func groupEvents(events: [LaunchEvent]) -> (today: [LaunchEvent], thisWeek: [LaunchEvent], future: [LaunchEvent], past: [LaunchEvent]) {
        var today: [LaunchEvent] = []
        var thisWeek: [LaunchEvent] = []
        var future: [LaunchEvent] = []
        var past: [LaunchEvent] = []
        
        let calendar = Calendar.current
        let now = Date()
        let startOfToday = calendar.startOfDay(for: now)
        let endOfToday = calendar.date(byAdding: .day, value: 1, to: startOfToday)!
        let startOfNextWeek = calendar.date(byAdding: .day, value: 7, to: startOfToday)!

        for event in events {
            if let eventDate = event.parsedDate {
                if eventDate >= startOfToday && eventDate < endOfToday {
                    today.append(event)
                } else if eventDate >= endOfToday && eventDate < startOfNextWeek {
                    thisWeek.append(event)
                } else if eventDate >= startOfNextWeek {
                    future.append(event)
                } else {
                    past.append(event)
                }
            } else {
                past.append(event)
            }
        }

        return (today, thisWeek, future, past)
    }
    
    @ViewBuilder
    private func EventRow(event: LaunchEvent) -> some View {
        NavigationLink(destination: LaunchEventDetailView(event: event)) {
            VStack(alignment: .leading, spacing: 4) {
                Text(event.title)
                    .font(.headline)

                Text(event.datetime)
                    .font(.subheadline)
                    .foregroundColor(.gray)

                if let missionID = event.mission_id {
                    Text("Mission ID: \(missionID)")
                        .font(.caption)
                        .foregroundColor(.blue)
                }

                if let providerIDString = event.spacedevs_id,
                   let providerID = Int(providerIDString),
                   let provider = providerFetcher.providers.first(where: { $0.spacedevs_id == providerID }) {
                    Text(provider.name)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
        }
    }
}
