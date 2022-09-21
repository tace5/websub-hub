package hub

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Topic struct {
	wg          *sync.WaitGroup
	subscribers []url.URL
}

func NewTopic() *Topic {
	t := Topic{
		wg:          new(sync.WaitGroup),
		subscribers: []url.URL{},
	}

	return &t
}

func (topic Topic) subscribe(callbackUrl url.URL) {
	topic.wg.Add(1)

	go func() {
		defer topic.wg.Done()
		topic.subscribers = append(topic.subscribers, callbackUrl)
	}()

	topic.wg.Wait()
}

func (topic Topic) publish(data []byte) {
	buffer := bytes.NewBuffer(data)

	for _, cbUrl := range topic.subscribers {
		_, err := http.Post(cbUrl.String(), "application/json", buffer)
		if err != nil {
			log.Print(err)
		}
	}
}
