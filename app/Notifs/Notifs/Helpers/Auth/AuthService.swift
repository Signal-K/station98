//
//  AuthService.swift
//  Station98
//
//  Created by Liam Arbuckle on 7/9/2025.
//

import Foundation
import Appwrite

@MainActor
class AuthService: ObservableObject {
    static let shared = AuthService()

    private let account: Account
    @Published var isLoggedIn = false
    @Published var errorMessage: String? = nil

    private init() {
        let client = Client()
            .setEndpoint("http://localhost:8020/v1")
            .setProject("station126")
            .setSelfSigned(true)

        self.account = Account(client)
    }

    func login(email: String, password: String) async {
        do {
            _ = try await account.createSession(userId: email, secret: password)
            isLoggedIn = true
        } catch {
            errorMessage = error.localizedDescription
            isLoggedIn = false
        }
    }

    func signup(email: String, password: String, name: String) async {
        do {
            _ = try await account.create(userId: ID.unique(), email: email, password: password, name: name)
            await login(email: email, password: password)
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func logout() async {
        do {
            _ = try await account.deleteSessions()
            isLoggedIn = false
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
