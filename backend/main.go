package main

import (
	// "flag"
	// "fmt"
	"log"
	"net/http"
	// "os"

	// "github.com/heyimcarlos/chat-app/backend/api"
	"github.com/heyimcarlos/chat-app/backend/internal"
)

func main() {
	// listenAddr := flag.String("listenaddr", os.Getenv("PORT"), "server listen address")

	// server := api.NewServer(*listenAddr)
	// fmt.Println("Server running on: ", *listenAddr)
	// log.Fatal(server.FiberRun())

	setupAPI()
}

func setupAPI() {

	manager := internal.NewManager()

	// gorilla websocket handler
	http.HandleFunc("/ws", manager.ServeWS)
	log.Println("Server running on: ", ":3001")

	// init server
	log.Fatal(http.ListenAndServe(":3001", nil))
}
