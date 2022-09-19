package hub

import (
	"encoding/json"
	"log"
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

func (hub Hub) subscribe(topicName string, callbackUrl string) {
	topic := hub.topics[topicName]
	topic.subscribe(callbackUrl)
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
