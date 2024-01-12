package service

import (
	"sync"
	"time"
)

const (
	maxRequestsPerWindow = 10 // Maximum requests allowed within the window
	windowDuration       = 5 * time.Second
)

var (
	requestsInWindow  = 0
	requestTimestamps = make([]time.Time, 0, maxRequestsPerWindow)
)

func allowRequest() bool {
	now := time.Now()
	var lock sync.Mutex
	lock.Lock()
	for i := 0; i < len(requestTimestamps)-1; i++ {
		if now.Sub(requestTimestamps[i]) > windowDuration {
			requestTimestamps = requestTimestamps[i+1:]
			requestsInWindow--
		} else {
			break
		}
	}

	// If the window is full, deny the request
	if requestsInWindow >= maxRequestsPerWindow {
		return false
	}

	requestsInWindow++
	requestTimestamps = append(requestTimestamps, now)
	lock.Unlock()
	return true
}
