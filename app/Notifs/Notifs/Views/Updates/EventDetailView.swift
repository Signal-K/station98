//
//  EventDetailView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI
#if canImport(WebKit)
import WebKit
#endif

struct EventDetailView: View {
    @StateObject private var fetcher = EventFetcher()

    var body: some View {
        ScrollView {
            if fetcher.events.isEmpty {
                ProgressView("Loading events...")
            } else {
                VStack(alignment: .leading, spacing: 20) {
                    ForEach(fetcher.events) { event in
                        VStack(alignment: .leading, spacing: 20) {
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
                                    if let videoId = extractYouTubeID(from: video.url) {
                                        YouTubePlayerView(videoID: videoId)
                                            .frame(height: 320)
                                            .cornerRadius(12)
                                            .padding(.bottom)
                                    } else {
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
                            }
                        }
                    }
                }
                .padding()
            }
        }
        .navigationTitle("Event Details")
        .task {
            await fetcher.fetchEventsWithVideosOrUpdates()
        }
    }
}

func extractYouTubeID(from url: String) -> String? {
    let pattern = #"(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([\w\-]+)"#
    if let regex = try? NSRegularExpression(pattern: pattern),
       let match = regex.firstMatch(in: url, range: NSRange(url.startIndex..., in: url)),
       let range = Range(match.range(at: 1), in: url) {
        return String(url[range])
    }
    return nil
}

struct YouTubePlayerView: UIViewRepresentable {
    let videoID: String

    func makeUIView(context: Context) -> WKWebView {
        return WKWebView()
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {
        let embedHTML = """
            <!DOCTYPE html>
            <html>
            <body style="margin:0">
            <iframe width="100%" height="200%" src="https://www.youtube.com/embed/\(videoID)?playsinline=1" frameborder="0" allowfullscreen></iframe>
            </body>
            </html>
            """
        uiView.loadHTMLString(embedHTML, baseURL: nil)
    }
}
