//
//  AstronautsView.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import SwiftUI

struct AstronautsView: View {
    @State private var astronauts: [Astronaut] = []
    @State private var isLoading = true
    @State private var errorMessage: String?

    var body: some View {
        NavigationView {
            ScrollView {
                if isLoading {
                    ProgressView("Fetching Astronauts...")
                        .padding()
                } else if let errorMessage = errorMessage {
                    Text("❌ Error: \(errorMessage)")
                        .foregroundColor(.red)
                        .padding()
                } else {
                    LazyVStack(spacing: 16) {
                        ForEach(astronauts) { astronaut in
                            AstronautCardView(astronaut: astronaut)
                                .padding(.horizontal)
                        }
                    }
                }
            }
            .navigationTitle("🧑‍🚀 Astronauts")
            .onAppear(perform: loadAstronauts)
        }
    }

    private func loadAstronauts() {
        isLoading = true
        AstronautService.shared.fetchAstronauts { result in
            DispatchQueue.main.async {
                isLoading = false
                switch result {
                case .success(let fetched):
                    astronauts = fetched
                case .failure(let error):
                    errorMessage = error.localizedDescription
                }
            }
        }
    }
}
