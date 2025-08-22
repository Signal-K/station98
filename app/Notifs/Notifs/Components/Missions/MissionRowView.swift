//
//  MissionRowView.swift
//  Station98
//
//  Created by Liam Arbuckle on 25/8/2025.
//

import SwiftUI

struct MissionRowView: View {
    let mission: Mission
    let events: [LaunchEvent]

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(mission.name)
                .font(.title3)
                .bold()

            if let description = mission.description {
                Text(description)
                    .font(.body)
            }

            if let orbit = mission.orbit {
                Text("Orbit: \(orbit)")
                    .font(.footnote)
                    .foregroundColor(.secondary)
            }

            if !events.isEmpty {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Upcoming Events")
                        .font(.subheadline)
                        .foregroundColor(.accentColor)
                        .padding(.top, 8)

                    ForEach(events) { event in
                        VStack(alignment: .leading, spacing: 2) {
                            Text(event.title)
                                .font(.body)
                                .bold()
                            Text(event.formattedDate)
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        .padding(.vertical, 4)
                    }
                }
            }
        }
        .padding(.vertical)
    }
}
