//
//  PadsGlobeView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI

struct PadsGlobeView: View {
    @StateObject private var fetcher = PadFetcher()
    @State private var isRotating = true
    @State private var showLabels = true
    @State private var selectedCountries: Set<String> = ["CHN", "GBR", "GUF", "IND", "JPN", "KAZ", "KOR", "NZL", "RUS", "USA"]
    
    var body: some View {
        ZStack {
            if fetcher.isLoading {
                ProgressView("Loading launchpads...")
            } else if let error = fetcher.error {
                Text("Error: \(error)")
                    .foregroundColor(.red)
            } else {
                EarthGlobeView(
                    pads: fetcher.pads,
                    isRotating: $isRotating,
                    showLabels: $showLabels,
                    selectedCountries: $selectedCountries
                )
                VStack {
                    Spacer()
                    EarthGlobeControls(
                        isRotating: $isRotating,
                        showLabels: $showLabels,
                        selectedCountries: $selectedCountries
                    )
                }
            }
        }
        .task {
            await fetcher.fetchPads()
        }
    }
}
