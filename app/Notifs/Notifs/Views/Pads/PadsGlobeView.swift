//
//  PadsGlobeView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI

struct PadsGlobeView: View {
    @StateObject private var fetcher = PadFetcher()
    
    var body: some View {
        Group {
            if fetcher.isLoading {
                ProgressView("Loading launchpads...")
            } else if let error = fetcher.error {
                Text("Error: \(error)")
                    .foregroundColor(.red)
            } else {
                EarthGlobeView(pads: fetcher.pads)
            }
        }
        
        .task {
            await fetcher.fetchPads()
        }
    }
}
