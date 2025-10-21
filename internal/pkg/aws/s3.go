package aws

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client   *s3.Client
	bucketName string
)

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("ACCESS_KEY"),
			os.Getenv("SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("AWS_BUCKET")

	log.Printf("AWS S3 initialized successfully (bucket: %s)", bucketName)
}

func UploadPrivateFile(ctx context.Context, fileHeader *multipart.FileHeader, folderName string) (string, error) {

	if s3Client == nil {
		InitS3()
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s/%d-%s", folderName, time.Now().Unix(), filepath.Base(fileHeader.Filename))

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	return key, err

}

func GetPresignedURL(ctx context.Context, key string, expiresIn time.Duration) (*string, error) {

	if s3Client == nil {
		InitS3()
	}

	presigner := s3.NewPresignClient(s3Client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return nil, err
	}

	return &req.URL, nil
}

func UpdateFile(ctx context.Context, oldKey string, newFile *multipart.FileHeader, folder string) (string, error) {
	
	if s3Client == nil {
		InitS3()
	}

	if oldKey != "" {
		if err := DeleteFile(ctx, oldKey); err != nil {
			log.Printf("⚠️ Failed to delete old file: %v", err)
		}
	}

	return UploadPrivateFile(ctx, newFile, folder)

}

func DeleteFile(ctx context.Context, key string) error {

	if s3Client == nil {
		InitS3()
	}
	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return nil
}

func GetS3Client() *s3.Client {
	if s3Client == nil {
		InitS3()
	}
	return s3Client
}

func GetBucketName() string {
	if bucketName == "" {
		bucketName = os.Getenv("AWS_BUCKET")
	}
	return bucketName
}
