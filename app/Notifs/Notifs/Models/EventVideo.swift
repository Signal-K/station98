//
//  EventVideo.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import Foundation

struct EventVideo: Identifiable, Codable {
    var id = UUID()
    let priority: Int
    let title: String
    let url: String
}
