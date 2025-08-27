//
//  YoutubeVideoPlayer.swift
//  Station98
//
//  Created by Liam Arbuckle on 27/8/2025.
//

import SwiftUI
#if canImport(WebKit)
import WebKit
#endif

struct YoutubePlayerView: UIViewRepresentable {
    let videoId: String
    
    func makeUIView(context: Context) -> WKWebView {
        let webView = WKWebView()
        webView.scrollView.isScrollEnabled = false
        webView.configuration.allowsInlineMediaPlayback = true
        return webView
    }
    
    func updateUIView(_ uiView: WKWebView, context: Context) {
        guard let url = URL(string: "https://www.youtube.com/embed/\(videoId)?playsinline-1") else { return }
        let request = URLRequest(url: url)
        uiView.load(request)
    }
    
    func extractYouTubeID(from urlString: String) -> String? {
        if let url = URLComponents(string: urlString) {
            if url.host?.contains("youtu.be") == true {
                return url.path.dropFirst().description
            } else if url.host?.contains("youtube.com") == true {
                return url.queryItems?.first(where: { $0.name == "v" })?.value
            }
        }
        
        return nil
    }
}
