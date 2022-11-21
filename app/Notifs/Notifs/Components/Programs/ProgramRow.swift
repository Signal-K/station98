//
//  ProgramRow.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import SwiftUI

struct ProgramCard: View {
    let program: Program

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(program.name)
                .font(.title3)
                .bold()
                .foregroundColor(.primary)

            if let description = program.description {
                Text(description)
                    .font(.body)
                    .lineLimit(4)
                    .foregroundColor(.secondary)
            }

            HStack {
                Label("Start", systemImage: "calendar")
                Text(program.start_date ?? "Unknown")

                Spacer()

                Label("Type", systemImage: "scope")
                Text(program.type)
            }
            .font(.caption)

            if let wiki = program.wiki_url, let url = URL(string: wiki) {
                Link("🌐 Wikipedia", destination: url)
                    .font(.caption)
                    .foregroundColor(.blue)
            }
        }
        .padding()
        .background(.ultraThinMaterial)
        .clipShape(RoundedRectangle(cornerRadius: 16))
        .shadow(radius: 3)
    }
}
