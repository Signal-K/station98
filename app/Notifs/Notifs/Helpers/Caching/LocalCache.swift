//
//  LocalCache.swift
//  Station98
//
//  Created by Liam Arbuckle on 28/8/2025.
//

import Foundation

enum LocalCache {
    static func save<T: Codable>(_ object: T, to filename: String) {
        let url = cacheURL(for: filename)
        do {
            let data = try JSONEncoder().encode(object)
            try data.write(to: url)
        } catch {
            print("❌ Failed to save cache for \(filename): \(error)")
        }
    }
    
    static func load<T: Codable>(from filename: String, as type: T.Type) -> T? {
        let url = cacheURL(for: filename)
        do {
            let data = try Data(contentsOf: url)
            return try JSONDecoder().decode(T.self, from: data)
        } catch {
            print("⚠️ Failed to load cache for \(filename): \(error)")
            return nil
        }
    }
    
    private static func cacheURL(for filename: String) -> URL {
        FileManager.default.urls(for: .cachesDirectory, in: .userDomainMask)[0]
            .appendingPathComponent(filename)
    }
}
