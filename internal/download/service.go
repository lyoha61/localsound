package download

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"localsound-backend/internal/song"
	"localsound-backend/internal/storage"
	"localsound-backend/internal/ytdlp"
)

type DownloadRequest struct {
	URL string `json:"url"`
}

type Service struct {
	ytdlpService   *ytdlp.Service
	storageService *storage.StorageService
	songService    *song.Service
	logger         *zap.SugaredLogger
}

func NewService(
	ytdlpService *ytdlp.Service,
	storageService *storage.StorageService,
	songService *song.Service,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		ytdlpService:   ytdlpService,
		storageService: storageService,
		songService:    songService,
		logger:         logger,
	}
}

func (s *Service) Download(ctx context.Context, url string, options ytdlp.DownloadOptions) (string, error) {
	outputDir := s.storageService.GetSongsDirectory()

	audioPath, err := s.ytdlpService.Download(ctx, url, options, outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to download via yt-dlp: %w", err)
	}

	if _, err := s.songService.Create(ctx, audioPath); err != nil {
		s.logger.Errorf("Failed to persist song metadata for %s: %v", audioPath, err)
		// продолжаем, так как аудиофайл уже скачан
	}

	s.logger.Infof("URL processed: %s", url)
	return audioPath, nil
}
