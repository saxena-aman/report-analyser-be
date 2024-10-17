package uploadserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Generic function to handle internal server errors
func InternalServerError(c *gin.Context, message string) error {
	if message == "" {
		message = "Internal server error"
	}

	// Use the generic JSONResponse function to send the response
	JSONResponse(c, http.StatusInternalServerError, message, nil)
	return errors.New(message) // Return the error for consistency
}

// Generic function to create JSON responses
func JSONResponse(c *gin.Context, statusCode int, message string, data map[string]interface{}) {
	response := gin.H{"message": message}
	for key, value := range data {
		response[key] = value
	}
	c.JSON(statusCode, response)
}
