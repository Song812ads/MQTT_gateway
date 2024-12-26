package main

import (
	"fmt"
	"log"

	"github.com/gateway/api"
)

func main() {
	server := api.NewAPIServer(":12345")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server stop")
}
