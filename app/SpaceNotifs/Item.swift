//
//  Item.swift
//  SpaceNotifs
//
//  Created by Liam Arbuckle on 15/8/2025.
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
