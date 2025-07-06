//
//  SpacewalksView.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import SwiftUI

struct SpacewalksView: View {
    @StateObject private var viewModel = SpacewalksViewModel()

    var body: some View {
        NavigationView {
            List {
                if viewModel.loading {
                    ProgressView("Loading spacewalks...")
                } else {
                    if !viewModel.recent.isEmpty {
                        Section(header: Text("Upcoming & Recent")) {
                            ForEach(viewModel.recent) { sw in
                                SpacewalkRowView(spacewalk: sw)
                            }
                        }
                    }

                    if !viewModel.archived.isEmpty {
                        Section(header: Text("Archived")) {
                            ForEach(viewModel.archived) { sw in
                                SpacewalkRowView(spacewalk: sw)
                            }
                        }
                    }
                }
            }
            .navigationTitle("Spacewalks")
        }
        .task {
            await viewModel.fetchSpacewalks()
        }
    }
}
