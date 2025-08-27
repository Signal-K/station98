//
//  ContentView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 16/8/2025.
//

import SwiftUI

struct ContentView: View {
    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            LaunchEventsView()
                .tabItem {
                    Label("Launches", systemImage: "calendar")
                }
                .tag(0)

            MissionListView()
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
}
