package local_storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	BasePath string
	logger   *zap.Logger
}

func NewService(basePath string, logger *zap.Logger) (*Service, error) {

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	logger.Info("LocalStorageService initialized", zap.String("BasePath", basePath))

	return &Service{
		BasePath: basePath,
		logger:   logger,
	}, nil
}

func (s *Service) UploadFile(subDirectory string, filename string, file io.Reader) (string, error) {
	s.logger.Info("Uploading file to local storage", zap.String("subDirectory", subDirectory), zap.String("filename", filename))

	id := uuid.New().String()[:8]
	now := time.Now()
	date := now.Format("20060102")
	timeStr := now.Format("1504")

	originalFilename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	fileExtension := filepath.Ext(filename)
	filename = fmt.Sprintf("%s-%s-%s-%s%s", originalFilename, date, timeStr, id, fileExtension)

	subDirectory = fmt.Sprintf("%s/%s", subDirectory, now.Format("2006-01-02"))

	fileDir := filepath.Join(s.BasePath, subDirectory)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create sub-directory: %w", err)
		}
	}

	filePath := filepath.Join(fileDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	relativePath := filepath.Join(subDirectory, filename)
	s.logger.Info("File uploaded successfully", zap.String("path", relativePath))

	return relativePath, nil
}

func (s *Service) DeleteFile(relativePath string) error {
	filePath := filepath.Join(s.BasePath, relativePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.logger.Warn("file not found, skipping deletion", zap.String("path", filePath))
		return nil // File doesn't exist, so no need to return an error
	}
	if err := os.Remove(filePath); err != nil {
		s.logger.Error("failed to delete file", zap.String("path", filePath), zap.Error(err))
		return fmt.Errorf("failed to delete file: %w", err)
	}
	s.logger.Info("successfully deleted file", zap.String("path", filePath))
	return nil
}
