package main

import (
	"flag"
	"fmt"
	"github.com/heyimcarlos/chat-app/backend/api"
	"log"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3001", "server listen address")

	server := api.NewServer(*listenAddr)
	fmt.Println("Server running on: ", *listenAddr)
	log.Fatal(server.Run())
}
