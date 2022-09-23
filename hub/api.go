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
	return &API{hub: hub}
}

func (api API) start() {
	http.HandleFunc("/", api.handleSubscriberAction)
	http.HandleFunc("/notify", api.handleNotifySubscribers)

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
	secret := r.FormValue("hub.secret")
	callbackUrl, err := url.Parse(r.FormValue("hub.callback"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.hub.subscriberAction(mode, topic, *callbackUrl, secret)
}

func (api API) handleNotifySubscribers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only the POST method is supported")
	}

	topic := r.FormValue("hub.topic")
	data := r.FormValue("data")

	api.hub.notifySubscribers(topic, data)
}
