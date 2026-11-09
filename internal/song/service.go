package song

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"localsound-backend/internal/lyrics"
	"localsound-backend/internal/songfile"
	"localsound-backend/internal/ytdlp"

	"go.uber.org/zap"
)

type Service struct {
	songRepo        *Repository
	songFileRepo    *songfile.Repository
	metadataService *MetadataService
	lyricsService   *lyrics.Service
	ytdlpService    *ytdlp.Service
	logger          *zap.SugaredLogger
	importQueue     chan ImportTask
}

type ImportTask struct {
	Format songfile.SongFormat
	Artist string
	Title  string
}

func NewService(
	songRepo *Repository,
	songFileRepo *songfile.Repository,
	metadataService *MetadataService,
	lyricsService *lyrics.Service,
	ytdlpService *ytdlp.Service,
	logger *zap.SugaredLogger,
) *Service {
	return &Service{
		songRepo:        songRepo,
		songFileRepo:    songFileRepo,
		metadataService: metadataService,
		lyricsService:   lyricsService,
		ytdlpService:    ytdlpService,
		logger:          logger,
		importQueue:     make(chan ImportTask, 1000),
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

func (s *Service) GetSongs(ctx context.Context) ([]SongDto, error) {
	rawSongs, err := s.songRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get songs from DB: %w", err)
	}

	songs := make([]SongDto, 0, len(rawSongs))

	for _, rawSong := range rawSongs {
		song := SongDto{
			Title:     rawSong.Title,
			Artist:    rawSong.Artist,
			Album:     rawSong.Album,
			Duration:  rawSong.Duration,
			CreatedAt: rawSong.CreatedAt,
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Service) QueueBatchImport(ctx context.Context, req ImportBatchRequest) error {
	for _, track := range req.Tracks {
		task := ImportTask{
			Artist: track.Artist,
			Title:  track.Title,
			Format: req.Format,
		}

		select {
		case s.importQueue <- task:
		case <-ctx.Done():
			return ctx.Err()
		default:
			return errors.New("import queue is full")
		}
	}
	return nil
}

func (s *Service) SearchSongs(ctx context.Context, query string) ([]ytdlp.SearchResultItem, error) {
	resulst, err := s.ytdlpService.SearchAll(ctx, query, 5)
	if err != nil {
		return nil, fmt.Errorf("Failed to search songs '%s': '%w'", query, err)
	}

	if resulst == nil {
		return []ytdlp.SearchResultItem{}, nil
	}
	return resulst, nil
}
