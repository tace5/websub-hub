package main

import (
	"net/url"
	"time"
)

type Subscriber struct {
	callback       url.URL
	secret         string
	expirationTime time.Time
}

func NewSubscriber(callback url.URL, secret string, leaseSeconds int) *Subscriber {
	return &Subscriber{
		callback:       callback,
		secret:         secret,
		expirationTime: time.Now().Add(time.Second * time.Duration(leaseSeconds)),
	}
}
