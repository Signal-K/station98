//
//  CalendarHelper.swift
//  Station98
//
//  Created by Liam Arbuckle on 28/8/2025.
//

import EventKit

class CalendarHelper {
    static let eventStore = EKEventStore()
    
    static func requestAccess(completion: @escaping (Bool) -> Void) {
        eventStore.requestAccess(to: .event) { granted, _ in
            DispatchQueue.main.async {
                completion(granted)
            }
        }
    }
    
    static func addLaunchToCalendar(event: LaunchEvent, provider: LaunchProvider?) {
        requestAccess { granted in
            guard granted else { return  }
            
            let ekEvent = EKEvent(eventStore: eventStore)
            ekEvent.calendar = eventStore.defaultCalendarForNewEvents
            
            let providerName = provider?.name ?? "Unknown Agency"
            ekEvent.title = "Launch: \(event.title) by \(providerName)"
            ekEvent.notes = event.description ?? ""
            
            if let video = event.vid_urls?.first {
                ekEvent.url = URL(string: video.url)
            }
            
            if let startDate = event.parsedDate {
                ekEvent.startDate = startDate
                ekEvent.endDate = startDate.addingTimeInterval(3600)
            } else {
                return
            }
            
            do {
                try eventStore.save(ekEvent, span: .thisEvent)
            } catch {
                print("Failed to save event: \(error)")
            }
        }
    }
}
