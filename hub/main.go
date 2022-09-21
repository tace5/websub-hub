package hub

func main() {
	hub := *NewHub()
	api := NewAPI(hub)
	api.start()
}
