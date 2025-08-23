//
//  ContentView.swift
//  Notifs
//
//  Created by Liam Arbuckle on 16/8/2025.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        TabView {
            LaunchEventsView()
                .tabItem {
                    Label("Launches", systemImage: "calendar")
                }

            LaunchProvidersView()
                .tabItem {
                    Label("Providers", systemImage: "globe")
                }
        }
    }
}
