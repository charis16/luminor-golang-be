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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

var R2Client *s3.Client
var R2BucketName string
var R2Endpoint string
var R2PublicURL string

func InitR2() {
	accessKey := GetEnvOrPanic("R2_ACCESS_KEY_ID")
	secretKey := GetEnvOrPanic("R2_SECRET_ACCESS_KEY")
	R2Endpoint = GetEnvOrPanic("R2_ENDPOINT")
	R2BucketName = GetEnvOrPanic("R2_BUCKET_NAME")
	R2PublicURL = GetEnvOrPanic("R2_PUBLIC_URL")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           R2Endpoint,
				SigningRegion: "auto",
			}, nil
		})),
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect to R2: %v", err)
	}

	R2Client = s3.NewFromConfig(cfg)
	log.Println("✅ R2 client initialized")
}

func UploadToR2(file multipart.File, fileHeader *multipart.FileHeader, prefix string) (string, error) {
	defer file.Close()

	bucket := GetEnvOrPanic("R2_BUCKET_NAME")
	publicURL := GetEnvOrPanic("R2_PUBLIC_URL")

	// 1. Clean nama file
	cleanFilename := strings.ReplaceAll(fileHeader.Filename, " ", "-")
	timestamp := time.Now().Format("20060102-150405")
	uniqueSuffix := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d_%s", timestamp, uniqueSuffix, cleanFilename)

	// 2. Gabungkan prefix kalau ada
	objectName := filename
	if prefix != "" {
		objectName = fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), filename)
	}

	// 3. Ambil Content-Type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 4. Upload ke R2
	_, err := R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// 5. Generate Public URL
	url := fmt.Sprintf("%s/%s", strings.TrimRight(publicURL, "/"), objectName)
	return url, nil
}

func StreamImageFromR2(c *gin.Context, filename string, contentType string, cacheDuration time.Duration) {
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	resp, err := R2Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(R2BucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image from storage"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", fmt.Sprintf("max-age=%.0f", cacheDuration.Seconds()))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}

func DeleteFromR2(bucket string, fileURL string) error {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	r2BaseURL := strings.TrimSuffix(GetEnvOrPanic("R2_PUBLIC_URL"), "/")

	// Pastikan URL berasal dari R2 bucket yang kita pakai
	if !strings.HasPrefix(fileURL, r2BaseURL) {
		return fmt.Errorf("URL does not match R2 base URL")
	}

	// Ambil path setelah base URL → jadi nama objeknya (bisa termasuk folder)
	objectKey := strings.TrimPrefix(parsed.Path, "/")

	// Pastikan tidak kosong
	if objectKey == "" {
		return fmt.Errorf("empty object key")
	}

	// Hapus dari R2
	_, err = R2Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(GetEnvOrPanic("R2_BUCKET_NAME")),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete R2 object: %w", err)
	}

	return nil
}
