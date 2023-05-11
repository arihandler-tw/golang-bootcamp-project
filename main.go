package main

import (
	"gin-exercise/pkg/product/broker"
	"gin-exercise/pkg/server"
	"os"
)

func main() {
	kc := os.Args[1]

	if kc == "consumer" {
		broker.Consumer()
		os.Exit(0)
	}

	router := server.SetupRoutes()
	router.Run("localhost:8080")
}
