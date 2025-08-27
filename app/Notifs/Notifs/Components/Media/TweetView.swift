//
//  TweetView.swift
//  Station98
//
//  Created by Liam Arbuckle on 27/8/2025.
//

import SwiftUI
import WebKit

struct TweetView: UIViewRepresentable {
    let url: URL

    func makeUIView(context: Context) -> WKWebView {
        let webView = WKWebView()
        webView.scrollView.isScrollEnabled = false
        webView.isOpaque = false
        webView.backgroundColor = .clear
        return webView
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {
        let tweetHTML = """
        <!DOCTYPE html>
        <html>
          <head>
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <style>
              body {
                margin: 0;
                padding: 0;
                background-color: transparent;
              }
            </style>
          </head>
          <body>
            <blockquote class="twitter-tweet">
              <a href="\(url.absoluteString)"></a>
            </blockquote>
            <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
          </body>
        </html>
        """
        uiView.loadHTMLString(tweetHTML, baseURL: nil)
    }
}
