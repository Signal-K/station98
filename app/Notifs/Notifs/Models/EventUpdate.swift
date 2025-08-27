//
//  EventUpdate.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import Foundation

struct EventUpdate: Identifiable, Codable {
    let id: String
    let title: String
    let description: String
    let created_at: String
    
    var createdAtFormatted: String {
        let formatter = ISO8601DateFormatter()
        if let date = formatter.date(from: created_at) {
            let displayFormatter = DateFormatter()
            displayFormatter.dateStyle = .medium
            displayFormatter.timeStyle = .short
            return displayFormatter.string(from: date)
        }
        
        return created_at
    }
}

struct UpdateEntry: Identifiable {
    var id: String
    var title: String
    var description: String
    var created_at: String
}
