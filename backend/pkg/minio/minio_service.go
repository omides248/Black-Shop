package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
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

func (s *Service) UploadFile(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64) (string, error) {
	s.logger.Info("Uploading file to MinIO", zap.String("bucket", bucketName), zap.String("object", objectName))

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

	_, err = s.client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", s.minioPublicURL, bucketName, objectName)
	s.logger.Info("File uploaded successfully", zap.String("url", fileURL))

	return fileURL, nil
}

func (s *Service) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	s.logger.Info("Deleting file from MinIO", zap.String("bucket", bucketName), zap.String("object", objectName))
	return s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *Service) DownloadFile(ctx context.Context, bucketName string, objectName string) (io.Reader, error) {
	s.logger.Info("Downloading file from MinIO", zap.String("bucket", bucketName), zap.String("object", objectName))
	object, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return object, nil
}
