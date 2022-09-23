package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const PORT = 8080

type API struct {
	hub Hub
}

func NewAPI(hub Hub) *API {
	api := API{hub: hub}

	return &api
}

func (api API) start() {
	http.HandleFunc("/", api.handleSubscriberAction)
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

	mode := r.FormValue("hub.mode")
	topic := r.FormValue("hub.topic")
	callbackUrl, err := url.Parse(r.FormValue("hub.callback"))
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

	topic := r.FormValue("hub.topic")
	data := r.FormValue("data")

	api.hub.publish(topic, data)
}
