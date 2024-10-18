package uploadserver

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// Handler to upload a file to S3
func uploadFileToS3(c *gin.Context, s3Client *s3.S3) {
	// Parse the file from the request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		InternalServerError(c, "File upload failed")
		return
	}
	defer file.Close()

	// Call the DAO function to upload the file to S3
	fileURL, err := UploadFileToS3(s3Client, file, header)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Respond with the file's public URL
	JSONResponse(c, http.StatusOK, "File uploaded successfully", map[string]interface{}{
		"file_url": fileURL,
	})
}

// Handler to list S3 objects
func listS3Objects(c *gin.Context, s3Client *s3.S3) {
	// Call the DAO function to get the list of objects
	objects, err := ListS3Objects(s3Client)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Respond with the list of objects
	JSONResponse(c, http.StatusOK, "Lookup successfully", map[string]interface{}{
		"items": objects,
	})
}

// UploadReportHandler handles the file upload and user info
func UploadReportHandler(c *gin.Context, s3Client *s3.S3, dynamoDBClient *dynamodb.DynamoDB) {
	// Parse form data and file
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		InternalServerError(c, "Unable to parse form data")
		return
	}

	// Get user details from the form
	name := c.PostForm("name")
	email := c.PostForm("email")
	gender := c.PostForm("gender")
	age := c.PostForm("age")

	// Get the file from the request
	header, err := c.FormFile("reportFile")
	if err != nil {
		InternalServerError(c, "Error retrieving the file")
		return
	}

	// Open the file for reading
	file, err := header.Open()
	if err != nil {
		InternalServerError(c, "Unable to open file")
		return
	}
	defer file.Close()

	// Call DAO function to handle the business logic
	err = HandleReportUpload(s3Client, dynamoDBClient, name, email, gender, age, file, header)
	if err != nil {
		InternalServerError(c, "Failed to upload report")
		return
	}

	// On success, send a response back
	JSONResponse(c, http.StatusOK, "Report uploaded successfully", nil)
}
