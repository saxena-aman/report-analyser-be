package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"report-analyser-be/uploadserver"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Reusable AWS S3 client
var s3Client *s3.S3
var dynamoDBClient *dynamodb.DynamoDB

// Init function for initializing AWS session and S3 client
func init() {
	// Check if running in development environment
	env := os.Getenv("ENV")
	if env == "development" {
		// Load the .env file only in development
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		fmt.Println(".env file loaded in development environment")
	}
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		log.Fatalf("AWS credentials or region not provided")
	}

	// Reusing AWS session with a defined HTTP client timeout for better performance
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Create S3 client once and reuse it
	s3Client = s3.New(sess)
	dynamoDBClient = dynamodb.New(sess)
}

func main() {
	r := gin.Default()

	// Add gzip compression to reduce response sizes
	r.Use(gzip.Gzip(gzip.BestSpeed))

	uploadserver.InitializeRoutes(r, s3Client, dynamoDBClient)
	// Start the Gin server
	fmt.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
