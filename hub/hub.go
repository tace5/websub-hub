package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Hub The type that represents the hub itself
type Hub struct {
	topics map[string]*Topic
	client http.Client
}

// NewHub Factory function for the hub
func NewHub() *Hub {
	h := Hub{topics: map[string]*Topic{}}

	return &h
}

// Registers a new topic
func (hub Hub) registerTopic(topicName string) {
	hub.topics[topicName] = NewTopic()
	print(hub.topics[topicName].subscribers)
}

// Subscribes or unsubscribes a subscriber from a topic (Includes validation and verification of intent)
func (hub Hub) subscriberAction(mode string, topicName string, callback url.URL, secret string) {
	if !hub.validateSubscription(mode, topicName) {
		return
	}
	if !verifyIntent(mode, topicName, callback) {
		return
	}

	if mode == "subscribe" {
		hub.subscribe(topicName, callback, secret)
	} else {
		hub.unsubscribe(topicName, callback)
	}
}

// Validates a subscription request
func (hub Hub) validateSubscription(mode string, topicName string) bool {
	failureReason := ""

	if mode != "subscribe" && mode != "unsubscribe" {
		failureReason = "Mode must be set to either 'subscribe' or 'unsubscribe"
	}
	if _, topicExists := hub.topics[topicName]; !topicExists {
		failureReason = "The topic does not exist"
	}

	if failureReason != "" {
		return false
	}

	return true
}

// Verifies the intent of a subscriber by sending a get request with the mode, topic and a challenge string
func verifyIntent(mode string, topic string, callback url.URL) bool {
	challenge := randomString(15)
	params := callback.Query()
	params.Add("hub.mode", mode)
	params.Add("hub.topic", topic)
	params.Add("hub.challenge", challenge)

	callback.RawQuery = params.Encode()
	resp, err := http.Get(callback.String())
	if err != nil {
		log.Print(err)
		return false
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Print(err)
	} else if resp.StatusCode == http.StatusOK && string(body) == challenge {
		return true
	}

	return false
}

// Subscribes a subscriber to a topic
func (hub Hub) subscribe(topicName string, callback url.URL, secret string) {
	topic := hub.topics[topicName]
	topic.subscribe(callback, secret)
}

// Unsubscribes a subscriber to a topic
func (hub Hub) unsubscribe(topicName string, callback url.URL) {
	topic := hub.topics[topicName]
	topic.unsubscribe(callback)
}

// Notifies all subscribers of a topic
func (hub Hub) notifySubscribers(topicName string, data string) {
	topic := hub.topics[topicName]
	payload := map[string]string{"hub.topic": topicName, "data": data}
	jsonData, err := json.Marshal(payload)

	if err != nil {
		log.Panic(err)
	}

	for callback, secret := range topic.subscribers {
		hub.notifySubscriber(callback, secret, jsonData, topicName)
	}
}

// Notifies a subscriber about changes to a topic
func (hub Hub) notifySubscriber(callback url.URL, secret string, data []byte, topicName string) {
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write(data)
	signature := hex.EncodeToString(hash.Sum(nil))
	buffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", callback.String(), buffer)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Link", fmt.Sprintf("http://hub:8080; rel=hub, %s; rel=self", topicName))
	req.Header.Add("X-Hub-Signature", "sha512="+signature)

	resp, err := hub.client.Do(req)
	if err != nil {
		log.Print(err)
	} else if resp.StatusCode == http.StatusGone {
		hub.unsubscribe(topicName, callback)
	}
}
