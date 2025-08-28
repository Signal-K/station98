//
//  MissionCard.swift
//  Station98
//
//  Created by Liam Arbuckle on 28/8/2025.
//

import SwiftUI

struct MissionCard: View {
    let mission: Mission
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text(mission.name)
                .font(.headline)
            
            if let orbit = mission.orbit, !orbit.isEmpty {
                Text("Orbit: \(orbit)")
                    .font(.subheadline)
                    .foregroundColor(.blue)
            }
            
            if let description = mission.description, !description.isEmpty {
                Text(description)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
        }
        
        .padding(.vertical, 8)
    }
}
