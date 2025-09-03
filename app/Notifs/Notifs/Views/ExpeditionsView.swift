//
//  ExpeditionsView.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import SwiftUI

struct ExpeditionsView: View {
    @State private var expeditions: [Expedition] = []
    @State private var isLoading = true

    var body: some View {
        // Section expeditions into "This Year" and "Archived"
        let currentYear = Calendar.current.component(.year, from: Date())
        let thisYearExpeditions = expeditions.filter {
            let year = Int($0.endDate.prefix(4)) ?? 0
            return year == currentYear
        }

        let archivedExpeditions = expeditions.filter {
            let year = Int($0.endDate.prefix(4)) ?? 0
            return year < currentYear
        }
        return NavigationView {
            List {
                if isLoading {
                    ProgressView()
                        .frame(maxWidth: .infinity, alignment: .center)
                } else {
                    if !thisYearExpeditions.isEmpty {
                        Section(header: Text("This Year")) {
                            ForEach(thisYearExpeditions) { expedition in
                                ExpeditionRow(expedition: expedition)
                            }
                        }
                    }
                    if !archivedExpeditions.isEmpty {
                        Section(header: Text("Archived")) {
                            ForEach(archivedExpeditions) { expedition in
                                ExpeditionRow(expedition: expedition)
                            }
                        }
                    }
                }
            }
            .navigationTitle("Expeditions")
            .onAppear {
                ExpeditionService.fetchExpeditions { result in
                    DispatchQueue.main.async {
                        switch result {
                        case .success(let data):
                            print("✅ Fetched \(data.count) expeditions")
                            for exp in data {
                                print("→ \(exp.name) at \(exp.stationName)")
                            }
                            self.expeditions = data
                        case .failure(let error):
                            print("❌ Failed to fetch expeditions: \(error.localizedDescription)")
                        }
                        self.isLoading = false
                    }
                }
            }
        }
    }
}
