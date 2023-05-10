package main

import (
	"gin-exercise/pkg/server"
)

func main() {
	router := server.SetupRoutes()
	router.Run("localhost:8080")
}
