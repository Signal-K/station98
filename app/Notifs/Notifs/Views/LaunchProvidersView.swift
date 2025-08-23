//
//  LaunchProvidersView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 24/8/2025.
//

import SwiftUI

struct LaunchProvidersView: View {
    @StateObject private var fetcher = LaunchProviderFetcher()
    @State private var sortMethod: SortMethod = .name
    @State private var selectedProvider: LaunchProvider?

    enum SortMethod: String, CaseIterable {
        case name = "Name"
        case country = "Country"
    }

    var body: some View {
        NavigationView {
            Group {
                if fetcher.isLoading {
                    ProgressView("Loading providers...")
                } else if let error = fetcher.error {
                    Text("Error: \(error)")
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                } else {
                    ScrollView {
                        VStack(alignment: .leading, spacing: 16) {
                            Picker("Sort by", selection: $sortMethod) {
                                ForEach(SortMethod.allCases, id: \.self) { method in
                                    Text(method.rawValue)
                                }
                            }
                            .pickerStyle(SegmentedPickerStyle())
                            .padding(.horizontal)

                            let grouped = Dictionary(grouping: fetcher.filteredProviders.sorted {
                                switch sortMethod {
                                case .name:
                                    return $0.name < $1.name
                                case .country:
                                    return $0.displayCountryCode < $1.displayCountryCode
                                }
                            }) { provider in
                                switch sortMethod {
                                case .name:
                                    return provider.nameInitial
                                case .country:
                                    return provider.displayCountryCode
                                }
                            }

                            ForEach(grouped.keys.sorted(), id: \.self) { section in
                                VStack(alignment: .leading, spacing: 8) {
                                    Text(section)
                                        .font(.title2)
                                        .bold()
                                        .padding(.horizontal)

                                    ForEach(grouped[section]!) { provider in
                                        Button {
                                            selectedProvider = provider
                                        } label: {
                                            LaunchProviderRowView(provider: provider)
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
            .navigationTitle("🚀 Launch Providers")
            .sheet(item: $selectedProvider) { provider in
                LaunchProviderDetailsView(provider: provider, allEvents: fetcher.allEvents)
            }
        }
        .task {
            await fetcher.fetch()
        }
    }
}
