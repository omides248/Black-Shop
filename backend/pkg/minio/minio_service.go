package minio

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type Service struct {
	client         *minio.Client
	logger         *zap.Logger
	minioPublicURL string
}

func NewService(endpoint, accessKey, secretKey string, minioPublicURL string, logger *zap.Logger) (*Service, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	logger.Info("Successfully connected to MinIO")

	return &Service{
		client:         client,
		logger:         logger.Named("minio_service"),
		minioPublicURL: minioPublicURL,
	}, nil
}

func (s *Service) DownloadFile(ctx context.Context, bucketName string, objectName string) (io.Reader, error) {
	s.logger.Info("Downloading file from MinIO", zap.String("bucket", bucketName), zap.String("object", objectName))
	object, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return object, nil
}

func (s *Service) UploadFile(ctx context.Context, bucketName, prefix, filename string, file io.Reader, fileSize int64, contentType string) (string, error) {
	s.logger.Info("Uploading file to MinIO", zap.String("bucket", bucketName), zap.String("filename", filename))

	objectKey := s.generateObjectKey(prefix, filename)

	s.logger.Info("Uploading file to MinIO", zap.String("bucket", bucketName), zap.String("object", objectKey))

	found, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !found {
		s.logger.Info("Bucket not found, creating a new one", zap.String("bucket", bucketName))
		if err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// TODO check expiration
	_, err = s.client.PutObject(ctx, bucketName, objectKey, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Info("File uploaded successfully", zap.String("objectKey", objectKey))

	return objectKey, nil
}

func (s *Service) DeleteFile(ctx context.Context, bucketName string, relativePath string) error {
	s.logger.Info("Deleting file from MinIO", zap.String("bucket", bucketName), zap.String("object", relativePath))
	return s.client.RemoveObject(ctx, bucketName, relativePath, minio.RemoveObjectOptions{})
}

func (s *Service) GetFileURL(bucketName, relativePath string) string {
	return fmt.Sprintf("%s/%s/%s", s.minioPublicURL, bucketName, relativePath)
}

func (s *Service) generateObjectKey(prefix, filename string) string {
	fileExtension := filepath.Ext(filename)
	uuidFilename := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	return fmt.Sprintf("%s/%s/%s", prefix, time.Now().Format("2006/01/02"), uuidFilename)
}
