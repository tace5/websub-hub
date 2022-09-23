package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Hub struct {
	topics map[string]*Topic
}

func NewHub() *Hub {
	h := Hub{topics: map[string]*Topic{}}

	return &h
}

func (hub Hub) registerTopic(topicName string) {
	hub.topics[topicName] = NewTopic()
	print(hub.topics[topicName].subscribers)
}

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
		// TODO: call unsubscribe
	}
}

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

func (hub Hub) subscribe(topicName string, callback url.URL, secret string) {
	topic := hub.topics[topicName]
	topic.subscribe(callback, secret)
}

func (hub Hub) notifySubscribers(topicName string, data string) {
	topic := hub.topics[topicName]

	payload := map[string]string{"hub.topic": topicName, "data": data}
	jsonData, err := json.Marshal(payload)

	if err != nil {
		log.Fatal(err)
	}

	topic.notifySubscribers(jsonData)
}
