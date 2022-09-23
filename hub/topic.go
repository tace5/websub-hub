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
	subscribers []url.URL
}

func (topic *Topic) subscribe(callbackUrl url.URL) {
	topic.mu.Lock()
	defer topic.mu.Unlock()
	topic.subscribers = append(topic.subscribers, callbackUrl)
}

func (topic *Topic) publish(data []byte) {
	buffer := bytes.NewBuffer(data)
	
	for _, cbUrl := range topic.subscribers {
		_, err := http.Post(cbUrl.String(), "application/json", buffer)
		if err != nil {
			log.Print(err)
		}
	}
}
