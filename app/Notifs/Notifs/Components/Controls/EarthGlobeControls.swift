//
//  EarthGlobeControls.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI

struct EarthGlobeControls: View {
    @Binding var isRotating: Bool
    @Binding var showLabels: Bool
    @Binding var selectedCountries: Set<String>
    
    let allCountries: [String] = ["CHN", "GBR", "IND", "JPN", "GUF", "KAZ", "KOR", "NZL", "RUS", "USA"]

    var body: some View {
        VStack(spacing: 16) {
            HStack(spacing: 12) {
                Button(action: {
                    isRotating.toggle()
                }) {
                    Label(isRotating ? "Pause" : "Play", systemImage: isRotating ? "pause.fill" : "play.fill")
                        .padding(8)
                        .background(Color.white.opacity(0.2))
                        .clipShape(Capsule())
                }

                Toggle("Labels", isOn: $showLabels)
                    .toggleStyle(SwitchToggleStyle(tint: .blue))
            }
            .padding(.horizontal)

            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    ForEach(allCountries, id: \.self) { country in
                        Button(action: {
                            if selectedCountries.contains(country) {
                                selectedCountries.remove(country)
                            } else {
                                selectedCountries.insert(country)
                            }
                        }) {
                            // Use .text modifier for correct accessibility and font size, and remove any extra labels.
                            Text(countryFlag(for: country))
                                .font(.system(size: 24))
                                .accessibilityHidden(true)
                                .padding(6)
                                .background(selectedCountries.contains(country) ? Color.blue.opacity(0.7) : Color.gray.opacity(0.4))
                                .clipShape(Circle())
                        }
                    }
                }
                .padding(.horizontal)
            }
        }
        .padding(.vertical)
        .background(BlurView(style: .systemMaterialDark))
        .cornerRadius(20)
        .padding()
    }

    func countryFlag(for code: String) -> String {
        let flags: [String: String] = [
            "CHN": "🇨🇳", "GBR": "🇬🇧", "IND": "🇮🇳", "JPN": "🇯🇵",
            "GUF": "🇬🇫", "KAZ": "🇰🇿", "KOR": "🇰🇷", "NZL": "🇳🇿",
            "RUS": "🇷🇺", "USA": "🇺🇸"
        ]
        return flags[code] ?? "❓"
    }
}

struct EarthGlobeControls_Previews: PreviewProvider {
    static var previews: some View {
        EarthGlobeControls(isRotating: .constant(true), showLabels: .constant(true), selectedCountries: .constant(["USA", "CHN"]))
            .background(Color.black)
    }
}
