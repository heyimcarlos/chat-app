package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/heyimcarlos/chat-app/backend/api"
)

func main() {
	listenAddr := flag.String("listenaddr", os.Getenv("PORT"), "server listen address")

	server := api.NewServer(*listenAddr)
	fmt.Println("Server running on: ", *listenAddr)
	log.Fatal(server.Run())
}
