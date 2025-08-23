//
//  LaunchProvidersRow.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import SwiftUI

struct LaunchProviderRowView: View {
    let provider: LaunchProvider

    var body: some View {
        HStack(alignment: .center, spacing: 12) {
            AsyncImage(url: URL(string: provider.image_url ?? "")) { phase in
                switch phase {
                case .empty:
                    ProgressView()
                        .frame(width: 48, height: 48)
                case .success(let image):
                    image
                        .resizable()
                        .scaledToFit()
                        .frame(width: 48, height: 48)
                        .cornerRadius(6)
                case .failure:
                    Image(systemName: "photo")
                        .frame(width: 48, height: 48)
                @unknown default:
                    EmptyView()
                }
            }

            VStack(alignment: .leading, spacing: 2) {
                Text(provider.name)
                    .font(.headline)
                Text(provider.displayCountryCode)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()
        }
        .padding(.horizontal)
    }
}
