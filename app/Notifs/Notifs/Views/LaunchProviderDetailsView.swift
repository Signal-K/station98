//
//  LaunchProviderDetailsView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import SwiftUI

struct LaunchProviderDetailsView: View {
    let provider: LaunchProvider
    let allEvents: [LaunchEvent]

    var filteredEvents: [LaunchEvent] {
        return allEvents.filter { event in
            if let eventID = event.spacedevs_id, let providerID = provider.spacedevs_id {
                return eventID == String(providerID)
            }
            return false
        }
    }

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                Text(provider.name)
                    .font(.title)
                    .bold()

                if let description = provider.description, !description.isEmpty {
                    Text(description)
                        .font(.body)
                } else {
                    Text("No description available.")
                        .font(.body)
                        .foregroundColor(.secondary)
                }

                if let type = provider.type, !type.isEmpty {
                    Text("Type: \(type)")
                        .font(.subheadline)
                        .foregroundColor(.gray)
                }

                if !filteredEvents.isEmpty {
                    Divider()
                    Text("Upcoming Launches")
                        .font(.headline)
                        .padding(.top)

                    ForEach(filteredEvents) { event in
                        VStack(alignment: .leading, spacing: 4) {
                            Text(event.title)
                                .font(.subheadline)
                            Text(event.formattedDate)
                                .font(.caption)
                                .foregroundColor(.gray)
                        }
                        .padding(.vertical, 4)
                    }
                }

                Spacer()
            }
            .padding()
        }
        .presentationDetents([.medium, .large])
    }
}
