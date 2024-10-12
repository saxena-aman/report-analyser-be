package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type S3Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

// Reusable AWS S3 client
var s3Client *s3.S3

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
}

// Handler to list S3 objects
func listS3Objects(c *gin.Context) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 bucket name not provided"})
		return
	}

	// List objects in the S3 bucket
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list S3 bucket objects: %v", err)})
		return
	}

	// Prepare a list of objects
	var objects []S3Object
	for _, item := range resp.Contents {
		objects = append(objects, S3Object{
			Key:  *item.Key,
			Size: *item.Size,
		})
	}

	// Return the list as JSON
	c.JSON(http.StatusOK, objects)
}

// Handler to upload file to S3
func uploadFileToS3(c *gin.Context) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 bucket name not provided"})
		return
	}

	// Parse the file from the request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	// Generate a unique file name based on the current timestamp
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"), // Set the ACL to public-read to make the file publicly accessible
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file to S3: %v", err)})
		return
	}

	// Generate the file's public URL
	region := os.Getenv("AWS_REGION")
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	// Return the file URL
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
}

func main() {
	r := gin.Default()

	// Add gzip compression to reduce response sizes
	r.Use(gzip.Gzip(gzip.BestSpeed))

	// Define the endpoint for listing S3 objects
	r.GET("/s3objects", listS3Objects)

	// Define the endpoint for file upload
	r.POST("/upload", uploadFileToS3)

	// Start the Gin server
	fmt.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
