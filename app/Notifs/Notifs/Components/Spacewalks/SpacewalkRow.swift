//
//  SpacewalkRow.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import SwiftUI

struct SpacewalkRowView: View {
    let spacewalk: Spacewalk
    
    var body: some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(spacewalk.name)
                .font(.headline)
            
            if let expedition = spacewalk.expeditionName {
                Text("Expedition: \(expedition)")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }

            if let date = ISO8601DateFormatter().date(from: spacewalk.startTime) {
                Text(formattedDate(date))
                    .font(.caption)
                    .foregroundColor(.gray)
            }
        }
        .padding(.vertical, 6)
    }
    
    func formattedDate(_ date: Date) -> String {
        let fmt = DateFormatter()
        fmt.dateStyle = .medium
        fmt.timeStyle = .short
        return fmt.string(from: date)
    }
}
