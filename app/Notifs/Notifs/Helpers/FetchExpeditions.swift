//
//  FetchExpeditions.swift
//  Station98
//
//  Created by Liam Arbuckle on 3/9/2025.
//

import Foundation

class ExpeditionService {
    static func fetchExpeditions(completion: @escaping (Result<[Expedition], Error>) -> Void) {
        guard let url = URL(string: "http://127.0.0.1:8080/api/collections/expeditions/records?expand=station&sort=-start_date&perPage=30") else {
            completion(.failure(URLError(.badURL)))
            return
        }

        URLSession.shared.dataTask(with: url) { data, _, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(URLError(.badServerResponse)))
                return
            }

            do {
                let result = try JSONDecoder().decode(ExpeditionResponse.self, from: data)
                completion(.success(result.items))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }

    private struct ExpeditionResponse: Codable {
        let items: [Expedition]
    }
}
