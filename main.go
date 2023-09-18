package main

import "Demerzel-Events/api"

func main() {
	srv := api.NewServer(5001, api.BuildRoutes())
	srv.Listen()
}
