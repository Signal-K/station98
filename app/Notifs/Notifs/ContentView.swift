//
//  ContentView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 16/8/2025.
//

import SwiftUI
import Clerk

struct ContentView: View {
    @Environment(\.clerk) private var clerk
    @State private var authIsPresented = false
    @State private var showUserButton = true
    @State private var selectedTab = 0

    var body: some View {
        VStack {
            if let _ = clerk.user {
                if showUserButton {
                    UserButton()
                        .frame(width: 36, height: 36)
                        .transition(.scale)
                        .onAppear {
                            DispatchQueue.main.asyncAfter(deadline: .now() + 3) {
                                withAnimation {
                                    showUserButton = false
                                }
                            }
                        }
                } else {
                    TabView(selection: $selectedTab) {
                        LaunchEventsView()
                            .tabItem {
                                Label("Launches", systemImage: "calendar")
                            }
                            .tag(0)

                        MissionsView()
                            .tabItem {
                                Label("Missions", systemImage: "lightbulb")
                            }
                            .tag(2)

                        EventDetailView()
                            .tabItem {
                                Label("Event", systemImage: "video.bubble.left")
                            }
                            .tag(4)

                        LaunchProvidersView()
                            .tabItem {
                                Label("Agencies", systemImage: "person")
                            }
                            .tag(5)

                        PadsGlobeView()
                            .tabItem {
                                Label("Pads", systemImage: "house")
                            }
                            .tag(3)
                    }
                    .onChange(of: selectedTab) { newTab in
                        if newTab == 3 {
                            UITabBar.appearance().barTintColor = .black
                            UITabBar.appearance().backgroundColor = .black
                        } else {
                            UITabBar.appearance().barTintColor = .systemBackground
                            UITabBar.appearance().backgroundColor = .systemBackground
                        }
                    }
                }
            } else {
                Button("Sign in") {
                    authIsPresented = true
                }
            }
        }
        .sheet(isPresented: $authIsPresented) {
            AuthView()
        }
    }
}
