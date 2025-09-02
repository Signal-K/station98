//
//  Astronaut.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import Foundation

struct Astronaut: Codable, Identifiable {
    let id: String
    let name: String
    let role: String?
    let priority: Int?
    let number: Int?
    let status: String?
    let in_space: Bool
    let eva_time_total: String?
    let space_time_total: String?
    let dob: String?
    let nationality: String?
    let first_flight: String?
    let last_flight: String?
    let agency: String?
    let agency_type: String?
    let bio: String?
    let wikipedia_url: String?
    let flights_count: Int?
    let landings_count: Int?
    let spacewalks_count: Int?
    let time_in_space: String?
    let eva_time: String?
    let date_of_death: String?
    let is_human: Bool
}
