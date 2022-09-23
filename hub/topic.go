package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

func (topic *Topic) notifySubscribers(data []byte) {
	for callback, secret := range topic.subscribers {
		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(data)
		payload := hex.EncodeToString(hash.Sum(nil))
		buffer := bytes.NewBuffer([]byte(payload))

		_, err := http.Post(callback.String(), "application/json", buffer)
		if err != nil {
			log.Print(err)
		}
	}
}
