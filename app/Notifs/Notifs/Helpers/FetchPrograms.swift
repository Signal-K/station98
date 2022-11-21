//
//  FetchPrograms.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import Foundation

class ProgramService {
    static let shared = ProgramService()
    private init() {}

    func fetchPrograms(completion: @escaping (Result<[Program], Error>) -> Void) {
        guard let url = URL(string: "http://127.0.0.1:8080/api/collections/programs/records?perPage=30&sort=-start_date") else {
            completion(.failure(NSError(domain: "Invalid URL", code: 0)))
            return
        }

        URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(NSError(domain: "No data", code: 0)))
                return
            }

            do {
                let decoded = try JSONDecoder().decode(PocketbaseProgramResponse.self, from: data)
                completion(.success(decoded.items))
            } catch {
                print("❌ Decoding error:", error)
                completion(.failure(error))
            }
        }.resume()
    }
}

// Wrapper to match Pocketbase response structure
struct PocketbaseProgramResponse: Codable {
    let items: [Program]
}
