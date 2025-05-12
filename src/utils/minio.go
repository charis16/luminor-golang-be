package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	protocol := "http"
	if GetEnvOrPanic("MINIO_USE_SSL") == "true" {
		protocol = "https"
	}
	minioEndpoint := GetEnvOrPanic("MINIO_ENDPOINT")

	defer file.Close()

	// 1. Check bucket
	exists, err := MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Printf("❌ Failed to check bucket existence: %v\n", err)
		return "", fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		log.Printf("ℹ️ Bucket '%s' not found, creating...\n", bucketName)
		err := MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("❌ Failed to create bucket: %v\n", err)
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// 2. Clean file name
	cleanFilename := strings.ReplaceAll(fileHeader.Filename, " ", "-")

	// 3. Add timestamp and unique suffix
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	uniqueSuffix := now.UnixNano()
	objectName := fmt.Sprintf("%s_%d_%s", timestamp, uniqueSuffix, cleanFilename)

	// 3. Log content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		log.Println("⚠️ Content-Type is empty, setting to application/octet-stream")
		contentType = "application/octet-stream"
	}

	// 4. Debug info
	log.Printf("📤 Uploading object: %s (size: %d bytes, type: %s)\n", objectName, fileHeader.Size, contentType)

	// 5. Upload
	_, err = MinioClient.PutObject(context.Background(), bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("❌ Failed to upload to MinIO: %v\n", err)
		return "", err
	}

	// 6. Return final public URL
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, minioEndpoint, bucketName, objectName)
	log.Printf("✅ Upload success: %s\n", url)
	return url, nil
}

func StreamImageFromMinio(c *gin.Context, bucket, filename string, contentType string, cacheDuration time.Duration) {
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-type", contentType)

	presignedURL, err := MinioClient.PresignedGetObject(context.Background(), bucket, filename, cacheDuration, reqParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate presigned URL"})
		return
	}

	resp, err := http.Get(presignedURL.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image from storage"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Cache-Control", fmt.Sprintf("max-age=%.0f", cacheDuration.Seconds()))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}

func DeleteFromMinio(bucketName, objectName string) error {
	objName := GetObjectNameFromURL(objectName)
	err := MinioClient.RemoveObject(context.Background(), bucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		// If it's a known error from MinIO
		if minioErr, ok := err.(minio.ErrorResponse); ok {
			if minioErr.Code == "NoSuchKey" {
				// Silently skip missing file
				return nil
			}
		}
		// Any other error
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}
