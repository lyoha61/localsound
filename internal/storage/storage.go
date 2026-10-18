package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type StorageService struct {
	songsDir string
}

func NewStorageService(dir string) (*StorageService, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for storage dir: %w", err)
	}

	if err := os.MkdirAll(absDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", absDir, err)
	}

	return &StorageService{songsDir: absDir}, nil
}

func (s *StorageService) GetSongsDirectory() string {
	return s.songsDir
}
