package main

import (
	"net/url"
	"sync"
)

// Topic Represents a topic
type Topic struct {
	mu          sync.Mutex
	subscribers map[url.URL]string
}

// NewTopic Factory function for the Topic type
func NewTopic() *Topic {
	return &Topic{
		mu:          sync.Mutex{},
		subscribers: map[url.URL]string{},
	}
}

// Adds a callback url and the provided secret to the list of subscribers
func (topic *Topic) subscribe(callback url.URL, secret string) {
	topic.mu.Lock()
	defer topic.mu.Unlock()
	topic.subscribers[callback] = secret
}

// Removes the given callback from the subscribers list
func (topic *Topic) unsubscribe(callback url.URL) {
	topic.mu.Lock()
	defer topic.mu.Unlock()
	delete(topic.subscribers, callback)
}
