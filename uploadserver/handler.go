package uploadserver

import (
	"net/http"

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
