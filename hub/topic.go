package main

import (
	"net/url"
	"sync"
)

type Topic struct {
	mu          sync.Mutex
	subscribers map[url.URL]string
}

func NewTopic() *Topic {
	return &Topic{
		mu:          sync.Mutex{},
		subscribers: map[url.URL]string{},
	}
}

func (topic *Topic) subscribe(callback url.URL, secret string) {
	topic.mu.Lock()
	defer topic.mu.Unlock()
	topic.subscribers[callback] = secret
}

func (topic *Topic) unsubscribe(callback url.URL) {
	topic.mu.Lock()
	defer topic.mu.Unlock()
	delete(topic.subscribers, callback)
}
