package hub

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

type Topic struct {
	wg          *sync.WaitGroup
	subscribers []string
}

func NewTopic() *Topic {
	t := Topic{
		wg:          new(sync.WaitGroup),
		subscribers: []string{},
	}

	return &t
}

func (topic Topic) subscribe(callbackUrl string) {
	topic.wg.Add(1)

	go func() {
		defer topic.wg.Done()
		topic.subscribers = append(topic.subscribers, callbackUrl)
	}()

	topic.wg.Wait()
}

func (topic Topic) publish(data []byte) {
	buffer := bytes.NewBuffer(data)

	for _, url := range topic.subscribers {
		_, err := http.Post(url, "application/json", buffer)
		if err != nil {
			log.Print(err)
		}
	}
}
