//
//  MissionListView.swift
//  Station98
//
//  Created by Liam Arbuckle on 25/8/2025.
//

import SwiftUI

struct MissionListView: View {
    @State private var missionEventPairs: [(Mission, [LaunchEvent])] = []
    @State private var isLoading = true

    var body: some View {
        NavigationView {
            Group {
                if isLoading {
                    ProgressView("Loading Missions…")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else if missionEventPairs.isEmpty {
                    VStack {
                        Text("No missions available.")
                            .font(.headline)
                            .padding()
                        Spacer()
                    }
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
                } else {
                    List {
                        ForEach(Array(missionEventPairs.enumerated()), id: \.element.0.id) { _, pair in
                            MissionRowView(mission: pair.0, events: pair.1)
                        }
                    }
                }
            }
            .navigationTitle("Missions")
        }
        .onAppear {
            MissionEventFetcher.shared.fetchMissionsAndLinkedEvents { pairs in
                self.missionEventPairs = pairs
                self.isLoading = false
            }
        }
    }
}
