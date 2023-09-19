package main

import (
	"demerzel-events/api"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to conver PORT to integer: %v", err))
	}
	log.Printf("Listening on port %d...\n", port)

	srv := api.NewServer(uint16(port), api.BuildRoutesHandler())
	srv.Listen()
}
