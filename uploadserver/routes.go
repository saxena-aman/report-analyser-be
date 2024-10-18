package uploadserver

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// InitializeRoutes sets up the routes for the application
func InitializeRoutes(router *gin.Engine, s3Client *s3.S3, dynamoDBClient *dynamodb.DynamoDB) {
	// Define your routes and pass S3 client to the handlers
	router.POST("/upload", func(c *gin.Context) {
		uploadFileToS3(c, s3Client)
	})
	router.GET("/s3objects", func(c *gin.Context) {
		listS3Objects(c, s3Client)
	})
	// You can add more routes here as needed
	// router.GET("/some-other-route", ...)
	router.POST("/upload-report", func(c *gin.Context) {
		UploadReportHandler(c, s3Client, dynamoDBClient)
	})
}
