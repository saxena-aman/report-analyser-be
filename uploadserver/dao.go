package uploadserver

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
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

func HandleReportUpload(s3Client *s3.S3, dynamoDBClient *dynamodb.DynamoDB, name, email, gender, age string, file multipart.File, header *multipart.FileHeader) error {
	// 1. Upload the file to S3
	fileURL, err := UploadFileToS3(s3Client, file, header)
	if err != nil {
		return err
	}

	// 2. Create a new UUID for the user
	userId := uuid.New().String()

	// 3. Prepare the DynamoDB item to insert
	item := map[string]*dynamodb.AttributeValue{
		"UserId": {
			S: aws.String(userId),
		},
		"ReportIndex": {
			S: aws.String(userId), // You can generate this dynamically based on business logic
		},
		"Name": {
			S: aws.String(name),
		},
		"Email": {
			S: aws.String(email),
		},
		"Gender": {
			S: aws.String(gender),
		},
		"Age": {
			N: aws.String(age),
		},
		"Reports": {
			L: []*dynamodb.AttributeValue{
				{
					M: map[string]*dynamodb.AttributeValue{
						"FileUrl": {
							S: aws.String(fileURL),
						},
						"FileName": {
							S: aws.String(header.Filename),
						},
						"UploadedAt": {
							S: aws.String(time.Now().Format(time.RFC3339)),
						},
					},
				},
			},
		},
		"CreatedAt": {
			S: aws.String(time.Now().Format(time.RFC3339)),
		},
	}

	// 4. Insert the item into DynamoDB
	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("report-analyser-users"), // Replace with your table name
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}
