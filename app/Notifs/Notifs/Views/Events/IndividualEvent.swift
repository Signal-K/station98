//
//  IndividualEvent.swift
//  Station98
//
//  Created by Liam Arbuckle on 27/8/2025.
//

import SwiftUI
import EventKit

struct LaunchEventDetailView: View {
    let event: LaunchEvent
    @StateObject private var providerFetcher = LaunchProviderFetcher()

    var body: some View {
        ScrollView {
            content
        }
        .navigationTitle(event.title)
        .task {
            await providerFetcher.fetch()
        }
    }

    private var content: some View {
        VStack(alignment: .leading, spacing: 16) {
            headerSection
//            descriptionSection
//            missionSection
//            sourceLinkSection
            calendarSection
            updatesSection
            videosSection
            Spacer()
        }
        .padding()
    }

    private var headerSection: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(event.title)
                .font(.largeTitle)
                .bold()
            Text(event.datetime)
                .font(.subheadline)
                .foregroundColor(.gray)
        }
    }

//    private var descriptionSection: some View {
//        Group {
//            if let description = event.description {
//                Text(description)
//                    .font(.body)
//            }
//        }
//    }
//
//    private var missionSection: some View {
//        Group {
//            if let missionID = event.mission_id {
//                Text("Mission ID: \(missionID)")
//                    .font(.caption)
//                    .foregroundColor(.blue)
//            }
//        }
//    }

    private var sourceLinkSection: some View {
        Group {
            if let sourceURL = event.source_url, let url = URL(string: sourceURL) {
                Link("Source", destination: url)
                    .font(.callout)
            }
        }
    }
    
    private var calendarSection: some View {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd HH:mm:ss.SSSZ"
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        let parsedDate = formatter.date(from: event.datetime)
        let matchedProvider = providerFetcher.providers.first(where: { $0.spacedevs_id == Int(event.spacedevs_id ?? "") })

        return VStack(alignment: .leading) {
            Text("Debugging Calendar Section")
            Text("Event spacedevs_id: \(event.spacedevs_id ?? "nil")")
            Text("Fetched providers count: \(providerFetcher.providers.count)")
            Text("Datetime raw value: \(event.datetime)")

            if let provider = matchedProvider {
                Text("Matched provider: \(provider.name)")
            } else {
                Text("No provider match")
            }

            if let parsedDate = parsedDate {
                Text("Datetime valid: \(parsedDate.formatted())")
            } else {
                Text("Datetime invalid")
            }

            // Calendar export button
            if let parsedDate = parsedDate {
                Button(action: {
                    let eventStore = EKEventStore()
                    eventStore.requestAccess(to: .event) { granted, error in
                        if granted {
                            let event = EKEvent(eventStore: eventStore)
                            event.title = self.event.title
                            event.startDate = parsedDate
                            event.endDate = parsedDate.addingTimeInterval(3600) // default 1-hour duration
                            event.calendar = eventStore.defaultCalendarForNewEvents
                            do {
                                try eventStore.save(event, span: .thisEvent)
                                print("Calendar event added")
                            } catch {
                                print("Failed to save event: \(error)")
                            }
                        } else {
                            print("Access denied or error: \(error?.localizedDescription ?? "unknown")")
                        }
                    }
                }) {
                    Text("Add to Calendar")
                        .foregroundColor(.blue)
                        .padding(.top, 8)
                }
            }
        }
    }

    private var updatesSection: some View {
        let updates = event.updates ?? []
        let updateEntries = updates.map {
            UpdateEntry(id: $0.id, title: $0.title, description: $0.description, created_at: $0.created_at)
        }
        let highlightKeywords = ["scrubbed", "success", "go for", "rescheduled", "delayed", "on hold", "tweaked"]
        let highlightedUpdates = updateEntries.filter { update in
            let lower = update.title.lowercased()
            return highlightKeywords.contains(where: { lower.contains($0) })
        }
        let newsUpdates = updateEntries.filter { update in
            let lower = update.title.lowercased()
            return !highlightKeywords.contains(where: { lower.contains($0) })
        }

        return Group {
            if !highlightedUpdates.isEmpty {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Updates")
                        .font(.headline)
                    ForEach(highlightedUpdates) { update in
                        updateRow(for: update)
                    }
                }
            }
            if !newsUpdates.isEmpty {
                VStack(alignment: .leading, spacing: 8) {
                    Text("News")
                        .font(.headline)
                    ForEach(newsUpdates) { update in
                        updateRow(for: update)
                    }
                }
            }
        }
    }

    private func updateRow(for update: UpdateEntry) -> some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(alignment: .top, spacing: 8) {
                Circle()
                    .fill(colorForUpdate(update.title))
                    .frame(width: 8, height: 8)
                    .padding(.top, 5)

                VStack(alignment: .leading, spacing: 4) {
                    HStack {
                        Text(update.title)
                            .font(.subheadline)
                            .bold()
                        Spacer()
                        Text(update.createdAtFormattedRelative)
                            .font(.caption2)
                            .foregroundColor(.secondary)
                    }
                    if let url = URL(string: update.description),
                       url.host?.contains("x.com") == true || url.host?.contains("twitter.com") == true,
                       url.path.contains("/status/") {
                        TweetView(url: url)
                            .frame(height: 480)
                            .clipShape(RoundedRectangle(cornerRadius: 12))
                            .shadow(radius: 4)
                    } else if let url = URL(string: update.description) {
                        Button(action: {
                            UIApplication.shared.open(url)
                        }) {
                            Text("Open Article")
                                .font(.caption)
                                .foregroundColor(.blue)
                                .underline()
                        }
                    } else {
                        Text(update.description)
                            .font(.caption)
                            .foregroundColor(.primary)
                    }
                }
            }
        }
        .padding(.vertical, 6)
    }

    private func colorForUpdate(_ title: String) -> Color {
        let lowercased = title.lowercased()
        if lowercased.contains("scrubbed") {
            return .red
        } else if lowercased.contains("success") || lowercased.contains("go for") {
            return .green
        } else if lowercased.contains("rescheduled") ||
                    lowercased.contains("delayed") ||
                    lowercased.contains("on hold") ||
                    lowercased.contains("tweaked") {
            return .yellow
        } else {
            return .blue
        }
    }

    private var videosSection: some View {
        Group {
            if let videos = event.vid_urls, !videos.isEmpty {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Videos")
                        .font(.headline)
                    ForEach(videos) { video in
                        if let url = URL(string: video.url) {
                            Link(video.title, destination: url)
                                .font(.callout)
                        }
                    }
                }
            }
        }
    }
}

extension UpdateEntry {
    var createdAtFormattedRelative: String {
        let isoFormatter = ISO8601DateFormatter()
        isoFormatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        
        let fallbackFormatter = ISO8601DateFormatter()
        fallbackFormatter.formatOptions = [.withInternetDateTime]

        let date: Date
        if let parsed = isoFormatter.date(from: created_at) {
            date = parsed
        } else if let parsed = fallbackFormatter.date(from: created_at) {
            date = parsed
        } else {
            return created_at
        }

        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .full
        return formatter.localizedString(for: date, relativeTo: Date())
    }
}
