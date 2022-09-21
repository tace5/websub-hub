package hub

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

func (hub Hub) createTopic(topicName string) {
	hub.topics[topicName] = NewTopic()
}

func (hub Hub) subscriberAction(mode string, topic string, callback url.URL) {
	if !hub.validateSubscription(mode, topic, callback) {
		return
	}
	if !verifyIntent(mode, topic, callback) {
		return
	}

	if mode == "subscribe" {
		hub.subscribe(topic, callback)
	} else {
		// TODO: call unsubscribe
	}
}

func (hub Hub) validateSubscription(mode string, topic string, callback url.URL) bool {
	failureReason := ""

	if mode != "subscribe" && mode != "unsubscribe" {
		failureReason = "Mode must be set to either 'subscribe' or 'unsubscribe"
	}
	if _, topicExists := hub.topics[topic]; topicExists {
		failureReason = "The topic does not exist"
	}
	if mode == "unsubscribe" && contains(hub.topics[topic].subscribers, callback) {
		failureReason = "The callback is not subscribed to the topic specified"
	}

	if failureReason != "" {
		return false
	}

	return true
}

func verifyIntent(mode string, topic string, callback url.URL) bool {
	challenge := randomString(15)
	params := callback.Query()
	params.Add("mode", mode)
	params.Add("topic", topic)
	params.Add("challenge", challenge)

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

func (hub Hub) subscribe(topicName string, callbackURL url.URL) {
	topic := hub.topics[topicName]
	topic.subscribe(callbackURL)
}

func (hub Hub) publish(topicName string, msg string) {
	topic := hub.topics[topicName]

	data := map[string]string{"topic": topicName, "msg": msg}
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	topic.publish(jsonData)
}
