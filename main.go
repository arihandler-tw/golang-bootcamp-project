package main

import (
	"fmt"
	"gin-exercise/pkg/product/broker"
	"gin-exercise/pkg/server"
	"os"
)

func main() {
	go broker.Consumer()

	router := server.SetupRoutes()
	err := router.Run("localhost:8080")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to run app: %v", err)
		os.Exit(1)
	}
}
