package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Handler to list S3 objects
func testAPI(c *gin.Context) {

	// Return the list as JSON
	c.JSON(http.StatusOK, "Hello World!")
}

func main() {
	r := gin.Default()

	// Add gzip compression to reduce response sizes
	r.Use(gzip.Gzip(gzip.BestSpeed))

	// Define the endpoint for listing S3 objects
	r.GET("/test", testAPI)

	// Start the Gin server
	fmt.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
