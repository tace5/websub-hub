package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const PORT = 3000

type API struct {
	hub Hub
}

func NewAPI(hub Hub) *API {
	api := API{hub: hub}

	return &api
}

func (api API) start() {
	http.HandleFunc("/", api.handleSubscriberAction)
	http.HandleFunc("/topics/new", api.handleCreateTopic)
	http.HandleFunc("/publish", api.handlePublish)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (api API) handleSubscriberAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mode := body["mode"]
	topic := body["topic"]
	callbackUrl, err := url.Parse(body["callback"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.hub.subscriberAction(mode, topic, *callbackUrl)
}

func (api API) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	api.hub.publish(body["topic"], body["msg"])
}

func (api API) handleCreateTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}
}
