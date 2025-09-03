//
//  ExpeditionRow.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import SwiftUI

struct ExpeditionRow: View {
    let expedition: Expedition

    var body: some View {
        HStack(alignment: .top, spacing: 12) {
            if let patch = expedition.patches, let url = URL(string: patch) {
                AsyncImage(url: url) { image in
                    image
                        .resizable()
                        .scaledToFit()
                } placeholder: {
                    ProgressView()
                }
                .frame(width: 60, height: 60)
                .clipShape(RoundedRectangle(cornerRadius: 8))
            } else {
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 60, height: 60)
                    .overlay(
                        Image(systemName: "globe.europe.africa")
                            .foregroundColor(.gray)
                    )
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(expedition.name)
                    .font(.headline)

                Text(expedition.stationName)
                    .font(.subheadline)
                    .foregroundColor(.secondary)

                Text(expedition.formattedDate)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()
        }
        .padding(.vertical, 6)
    }
}
