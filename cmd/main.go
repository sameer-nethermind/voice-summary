package main

import (
	"fmt"
	"log"

	"github.com/cyberhawk12121/Saarthi/internal/api"
	"github.com/cyberhawk12121/Saarthi/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect to the DB
	conn, err := db.Create()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to the database")

	// 2. Initialize Gin
	router := gin.Default()

	// 3. Setup routes
	api.SetupRoutes(router, conn)

	// 4. Run server
	router.Run(":8080")
}
