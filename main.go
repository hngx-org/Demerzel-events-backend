package main

import "demerzel-events/api"

func main() {
	srv := api.NewServer(5001, api.BuildRoutes())
	srv.Listen()
}
