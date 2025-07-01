package utils

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	endpoint := "minio:9000"
	accessKeyID := "minio"
	secretAccessKey := "minio123"
	useSSL := false

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	log.Println("MinIO client initialized")

	// Retry logic: Wait for MinIO and bucket to be ready
	retries := 10
	for i := 1; i <= retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		exists, err := MinioClient.BucketExists(ctx, "documents")
		if err == nil {
			if !exists {
				err = MinioClient.MakeBucket(ctx, "documents", minio.MakeBucketOptions{})
				if err != nil {
					log.Fatalf("Failed to create bucket: %v", err)
				}
				log.Println("Bucket created: documents")
			} else {
				log.Println("Bucket already exists: documents")
			}
			break
		}

		log.Printf("Attempt %d: MinIO not ready yet (%v). Retrying in 2s...", i, err)
		time.Sleep(2 * time.Second)

		if i == retries {
			log.Fatalf("MinIO not reachable after %d retries: %v", retries, err)
		}
	}
}
