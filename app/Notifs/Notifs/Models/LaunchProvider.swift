//
//  LaunchProvider.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import Foundation

struct LaunchProvider: Codable, Identifiable {
    let id: String
    let name: String
    let abbrev: String?
    let country_code: String?
    let type: String?
    let logo_url: String?
    let image_url: String?
    let wiki_url: String?
    let info_url: String?
    let spacedevs_id: Int?
    
    let description: String?

    // Computed properties for grouping/display
    var nameInitial: String {
        return String(name.prefix(1)).uppercased()
    }

    var displayCountryCode: String {
        if let code = country_code {
            return String(code.prefix(3)).uppercased()
        }
        return "UNK"
    }
}

extension Array where Element == LaunchProvider {
    func provider(withID id: String) -> LaunchProvider? {
        return first(where: { $0.id == id })
    }
}
