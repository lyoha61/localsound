package song

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"
	"localsound-backend/internal/lyrics"
	"localsound-backend/internal/songfile"
)

type Service struct {
	songRepo        *Repository
	songFileRepo    *songfile.Repository
	metadataService *MetadataService
	lyricsService   *lyrics.Service
	logger          *zap.SugaredLogger
}

func NewService(
	songRepo *Repository,
	songFileRepo *songfile.Repository,
	metadataService *MetadataService,
	lyricsService *lyrics.Service,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		songRepo:        songRepo,
		songFileRepo:    songFileRepo,
		metadataService: metadataService,
		lyricsService:   lyricsService,
		logger:          logger,
	}
}

func (s *Service) Create(ctx context.Context, audioPath string) (*SongEntity, error) {
	meta, err := s.metadataService.GetData(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from %s: %w", audioPath, err)
	}

	format := songfile.ParseSongFormat(meta.Format)

	var artistPtr *string
	if meta.Artist != "" {
		artistPtr = &meta.Artist
	}
	var albumPtr *string
	if meta.Album != "" {
		albumPtr = &meta.Album
	}

	songEntity := &SongEntity{
		Title:    meta.Title,
		Artist:   artistPtr,
		Album:    albumPtr,
		Duration: meta.Duration,
	}

	if err := s.songRepo.Create(ctx, songEntity); err != nil {
		return nil, fmt.Errorf("failed to save song entity: %w", err)
	}

	songFileEntity := &songfile.SongFileEntity{
		SongID:   songEntity.ID,
		FilePath: audioPath,
		Format:   format,
		FileSize: meta.FileSize,
	}

	if err := s.songFileRepo.Create(ctx, songFileEntity); err != nil {
		return nil, fmt.Errorf("failed to save song file entity: %w", err)
	}

	if err := s.lyricsService.Create(ctx, songEntity.ID); err != nil {
		s.logger.Errorf("Failed to queue lyrics creation for song %s: %v", songEntity.ID, err)
	}

	s.logger.Infof("Audio file %s saved in DB with song id=%s", filepath.Base(audioPath), songEntity.ID)
	return songEntity, nil
}
