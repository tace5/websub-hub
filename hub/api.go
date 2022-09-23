package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Port The port to serve content to
const Port = 8080

// API The type that represents the http endpoints used to interact with the hub
type API struct {
	hub Hub
}

// NewAPI Factory function for the API type
func NewAPI(hub Hub) *API {
	return &API{hub: hub}
}

// Starts the http server
func (api API) start() {
	http.HandleFunc("/", api.handleSubscriberAction)
	http.HandleFunc("/notify", api.handleNotifySubscribers)

	err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Handler for when a subscribe/unsubscribe request comes in
func (api API) handleSubscriberAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}

	mode := r.FormValue("hub.mode")
	topic := r.FormValue("hub.topic")
	secret := r.FormValue("hub.secret")
	callbackUrl, err := url.Parse(r.FormValue("hub.callback"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.hub.subscriberAction(mode, topic, *callbackUrl, secret)
}

// Handler for when a publisher wants to notify subscribers of updates to a topic
func (api API) handleNotifySubscribers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}

	topic := r.FormValue("hub.topic")
	data := r.FormValue("data")

	api.hub.notifySubscribers(topic, data)
}
