//
//  Program.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import Foundation

struct Program: Identifiable, Codable {
    let id: String
    let api_id: Int
    let name: String
    let type: String
    let description: String?
    let start_date: String?
    let end_date: String?
    let info_url: String?
    let wiki_url: String?
    let image_url: String?
    let image_thumb_url: String?
    let api_url: String?
    let created: String?
    let updated: String?
}
