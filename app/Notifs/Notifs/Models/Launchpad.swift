//
//  Launchpad.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

struct Pad: Identifiable, Codable {
    let id: String
    let name: String
    let latitude: Double
    let longitude: Double
    let location_name: String?
    let country_code: String?
}

struct PadResponse: Codable {
    let items: [Pad]
}
