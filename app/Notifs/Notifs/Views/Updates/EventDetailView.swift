//
//  EventDetailView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI
import Foundation

struct EventDetailView: View {
    @StateObject private var fetcher = EventFetcher()

    var body: some View {
        NavigationView {
            content
                .navigationTitle("🚀 Updates")
        }
        .task {
            await fetcher.fetchEvents()
        }
    }

    private var content: some View {
        ScrollView {
            if fetcher.isLoading {
                ProgressView("Loading updates...")
            } else if let error = fetcher.error {
                Text("Error: \(error)")
                    .foregroundColor(.red)
                    .multilineTextAlignment(.center)
            } else {
                let feed = sortedUpdateFeed(from: fetcher.events)
                VStack(alignment: .leading, spacing: 12) {
                    ForEach(feed) { update in
                        updateRow(for: update)
                    }
                }
                .padding()
            }
        }
    }

    private func sortedUpdateFeed(from events: [LaunchEvent]) -> [AnyUpdate] {
        var allUpdates: [AnyUpdate] = []

        for event in events {
            let updates = (event.updates ?? []).compactMap { update -> AnyUpdate? in
                guard let date = parseISODate(update.created_at) else { return nil }
                return AnyUpdate(
                    id: update.id,
                    title: update.title,
                    description: update.description,
                    date: date,
                    type: .update
                )
            }

            let videoUpdates: [AnyUpdate] = (event.vid_urls ?? []).compactMap { video in
                guard let eventDate = event.parsedDate else { return nil }
                return AnyUpdate(
                    id: video.url,
                    title: video.title,
                    description: video.url,
                    date: eventDate,
                    type: .video
                )
            }

            allUpdates.append(contentsOf: updates + videoUpdates)
        }

        return allUpdates.sorted(by: { $0.date > $1.date })
    }

    private func parseISODate(_ string: String) -> Date? {
        let primary = ISO8601DateFormatter()
        primary.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = primary.date(from: string) {
            return date
        }
        let fallback = ISO8601DateFormatter()
        fallback.formatOptions = [.withInternetDateTime]
        return fallback.date(from: string)
    }

    private func updateRow(for update: AnyUpdate) -> some View {
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
                        Text(relativeTimeString(for: update.date))
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
                    } else if let url = URL(string: update.description),
                              update.type == .video {
                        Button(action: {
                            UIApplication.shared.open(url)
                        }) {
                            Text("Watch Video")
                                .font(.caption)
                                .foregroundColor(.blue)
                                .underline()
                        }
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

    private func relativeTimeString(for date: Date) -> String {
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .full
        return formatter.localizedString(for: date, relativeTo: Date())
    }

    struct AnyUpdate: Identifiable {
        enum UpdateType {
            case update
            case video
        }

        let id: String
        let title: String
        let description: String
        let date: Date
        let type: UpdateType
    }
}

struct VideoUpdateWrapper: Identifiable {
    let id: String
    let title: String
    let url: String
    let created_at: String
    let date: Date
}
