//
//  EventDetailView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI

struct EventDetailView: View {
    @StateObject private var fetcher = EventFetcher()

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                if let event = fetcher.events.first {
                    Text(event.title)
                        .font(.largeTitle)
                        .bold()
                        .padding(.top)

                    if let updates = event.updates, !updates.isEmpty {
                        Text("Recent Updates")
                            .font(.headline)

                        ForEach(updates) { update in
                            VStack(alignment: .leading, spacing: 6) {
                                Text(update.title)
                                    .font(.subheadline)
                                    .fontWeight(.semibold)
                                Text(update.description)
                                    .font(.footnote)
                                    .foregroundColor(.secondary)
                                Text(update.createdAtFormatted)
                                    .font(.caption2)
                                    .foregroundColor(.gray)
                            }
                            .padding()
                            .background(Color(.secondarySystemBackground))
                            .cornerRadius(10)
                        }
                    }

                    if let videoURLs = event.vid_urls, !videoURLs.isEmpty {
                        Text("Video Links")
                            .font(.headline)

                        ForEach(videoURLs) { video in
                            Link(destination: URL(string: video.url)!) {
                                HStack {
                                    Image(systemName: "play.rectangle.fill")
                                        .foregroundColor(.blue)
                                    Text(video.title)
                                    Spacer()
                                }
                                .padding()
                                .background(Color(.tertiarySystemBackground))
                                .cornerRadius(8)
                            }
                        }
                    }
                } else {
                    ProgressView("Loading event...")
                }
            }
            .padding()
        }
        .navigationTitle("Event Details")
        .task {
            await fetcher.fetchEventsWithVideosOrUpdates()
        }
    }
}
