//
//  Mission.swift
//  Station98
//
//  Created by Liam Arbuckle on 25/8/2025.
//

import Foundation

struct Mission: Identifiable, Codable {
    let id: String
    let name: String
    let description: String?
    let orbit: String?
}
