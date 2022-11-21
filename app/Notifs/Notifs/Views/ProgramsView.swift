//
//  ProgramsView.swift
//  Station98
//
//  Created by Liam Arbuckle on 2/9/2025.
//

import SwiftUI

struct ProgramsView: View {
    @State private var programs: [Program] = []
    @State private var isLoading = true
    @State private var errorMessage: String?

    var body: some View {
        NavigationView {
            Group {
                if isLoading {
                    ProgressView("Loading Programs…")
                } else if let error = errorMessage {
                    Text("Error: \(error)")
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                        .padding()
                } else {
                    ScrollView {
                        LazyVStack(spacing: 16) {
                            ForEach(programs) { program in
                                ProgramCard(program: program)
                                    .padding(.horizontal)
                            }
                        }
                        .padding(.top)
                    }
                }
            }
            .navigationTitle("🚀 Space Programs")
            .onAppear {
                fetch()
            }
        }
    }

    func fetch() {
        ProgramService.shared.fetchPrograms { result in
            DispatchQueue.main.async {
                self.isLoading = false
                switch result {
                case .success(let items):
                    self.programs = items
                case .failure(let error):
                    self.errorMessage = error.localizedDescription
                }
            }
        }
    }
}
