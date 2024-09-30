package main

import (
    "sync"
)

var (
    shortlinkStore = make(map[string]string)
    storeMutex     sync.Mutex
)

func generateShortlink(content string) string {
    // Simple shortlink generation logic (for demo purposes)
    return "short_" + content[:5]
}

func storeShortlink(shortlink, content string) {
    storeMutex.Lock()
    defer storeMutex.Unlock()
    shortlinkStore[shortlink] = content
}

func getShortlinkContent(shortlink string) (string, bool) {
    storeMutex.Lock()
    defer storeMutex.Unlock()
    content, found := shortlinkStore[shortlink]
    return content, found
}
