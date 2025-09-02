//
//  FetchAstronauts.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import Foundation

class AstronautService {
    static let shared = AstronautService()
    
    func fetchAstronauts(completion: @escaping (Result<[Astronaut], Error>) -> Void) {
        guard let url = URL(string: "http://localhost:8080/api/collections/astronauts/records?perPage=100") else {
            completion(.failure(URLError(.badURL)))
            return
        }

        var request = URLRequest(url: url)
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(URLError(.badServerResponse)))
                return
            }

            do {
                let decoded = try JSONDecoder().decode(PocketbaseResponseAstronaut.self, from: data)
                completion(.success(decoded.items))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
}

struct PocketbaseResponseAstronaut: Codable {
    let items: [Astronaut]
}
