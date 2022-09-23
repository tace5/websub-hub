package main

func main() {
	hub := NewHub()
	hub.registerTopic("/a/topic")
	api := NewAPI(*hub)
	api.start()
}
