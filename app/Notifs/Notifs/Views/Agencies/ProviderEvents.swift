//
//  ProviderEvents.swift
//  Station98
//
//  Created by Liam Arbuckle on 27/8/2025.
//

import SwiftUI

struct ProviderEventsView: View {
    let provider: LaunchProvider
    let allEvents: [LaunchEvent]

    var matchingEvents: [LaunchEvent] {
        allEvents.filter { $0.spacedevs_id == provider.name }
    }

    var body: some View {
        List {
            if matchingEvents.isEmpty {
                Text("No upcoming launches for this provider.")
                    .foregroundColor(.gray)
            } else {
                ForEach(matchingEvents) { event in
                    VStack(alignment: .leading, spacing: 4) {
                        Text(event.title)
                            .font(.headline)
                        Text(event.formattedDate)
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        if let mission = event.expand?.mission?.name {
                            Text("Mission: \(mission)")
                                .font(.caption)
                                .foregroundColor(.blue)
                        }
                    }
                    .padding(.vertical, 6)
                }
            }
        }
        .navigationTitle(provider.name)
    }
}
