package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	endpoint := GetEnvOrPanic("MINIO_ENDPOINT")
	accessKey := GetEnvOrPanic("MINIO_ACCESS_KEY")
	secretKey := GetEnvOrPanic("MINIO_SECRET_KEY")
	useSSL := GetEnvOrPanic("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MinIO: %v", err)
	}

	MinioClient = client
	log.Println("✅ MinIO client initialized")
}

func UploadToMinio(bucketName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	// Make sure bucket exists
	exists, err := MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err := MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	contentType := fileHeader.Header.Get("Content-Type")

	_, err = MinioClient.PutObject(context.Background(), bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), bucketName, objectName), nil
}
