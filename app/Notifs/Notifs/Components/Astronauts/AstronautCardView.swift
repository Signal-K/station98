//
//  AstronautCardView.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import SwiftUI

struct AstronautCardView: View {
    let astronaut: Astronaut

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(astronaut.name)
                .font(.title2)
                .fontWeight(.bold)

            if let role = astronaut.role {
                Text(role)
                    .font(.subheadline)
                    .foregroundColor(.gray)
            }

            if let agency = astronaut.agency {
                HStack {
                    Label(agency, systemImage: "building.2")
                        .font(.footnote)
                        .foregroundColor(.secondary)
                }
            }

            if let nationality = astronaut.nationality {
                Text("🇺🇳 Nationality: \(nationality)")
                    .font(.footnote)
            }

            if let dob = astronaut.dob {
                Text("🎂 DOB: \(dob)")
                    .font(.footnote)
            }

            if let flights = astronaut.flights_count {
                Text("🚀 Flights: \(flights)")
                    .font(.footnote)
            }

            if let bio = astronaut.bio {
                Text(bio)
                    .font(.footnote)
                    .lineLimit(4)
                    .foregroundColor(.secondary)
            }
        }
        .padding()
        .background(Color(.systemBackground).opacity(0.9))
        .cornerRadius(12)
        .shadow(radius: 4)
    }
}
