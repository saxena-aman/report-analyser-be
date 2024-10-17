package uploadserver

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadFileToS3 uploads a file to the specified S3 bucket and returns the file's URL
func UploadFileToS3(s3Client *s3.S3, file multipart.File, header *multipart.FileHeader) (string, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("S3 bucket name not provided")
	}

	// Generate a unique file name based on the current timestamp
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))

	// Upload the file to S3
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"), // Set the ACL to public-read to make the file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Generate the file's public URL
	region := os.Getenv("AWS_REGION")
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	return fileURL, nil
}

// ListS3Objects retrieves the list of objects from the S3 bucket
func ListS3Objects(s3Client *s3.S3) ([]S3Object, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("S3 bucket name not provided")
	}

	// List objects in the S3 bucket
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 bucket objects: %v", err)
	}

	// Prepare a list of objects
	var objects []S3Object
	for _, item := range resp.Contents {
		objects = append(objects, S3Object{
			Key:  *item.Key,
			Size: *item.Size,
		})
	}

	return objects, nil
}
