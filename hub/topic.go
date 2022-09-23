package main

import (
	"bytes"
	"log"
	"net/http"
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

func (topic *Topic) publish(data []byte) {
	buffer := bytes.NewBuffer(data)

	for callback, _ := range topic.subscribers {
		_, err := http.Post(callback.String(), "application/json", buffer)
		if err != nil {
			log.Print(err)
		}
	}
}
