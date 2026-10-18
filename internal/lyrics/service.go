package lyrics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo     *Repository
	provider Provider
	logger   *zap.SugaredLogger
}

func NewService(repo *Repository, provider Provider, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:     repo,
		provider: provider,
		logger:   logger,
	}
}

func (s *Service) Create(ctx context.Context, songID uuid.UUID) error {
	entity := &LyricsEntity{
		SongID:       songID,
		LyricsStatus: StatusPending,
		SyncStatus:   "PENDING",
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return fmt.Errorf("failed to create lyrics record: %w", err)
	}

	s.logger.Infof("Lyrics saved for song id=%s", songID)
	return nil
}

func (s *Service) ProcessPending(ctx context.Context) error {
	pending, err := s.repo.FindByStatus(ctx, StatusPending)
	if err != nil {
		return fmt.Errorf("failed to find pending lyrics: %w", err)
	}

	if len(pending) == 0 {
		return nil
	}

	s.logger.Infof("Processing %d pending lyrics", len(pending))

	for _, item := range pending {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second): // Небольшая пауза между запросами к внешнему API
		}

		artist := ""
		if item.SongArtist != nil {
			artist = *item.SongArtist
		}
		album := ""
		if item.SongAlbum != nil {
			album = *item.SongAlbum
		}

		s.logger.Infof("Searching lyrics for song: title='%s', artist='%s'", item.SongTitle, artist)

		req := LyricsRequest{
			TrackName:  item.SongTitle,
			ArtistName: artist,
			AlbumName:  album,
			Duration:   item.SongDuration,
		}

		lyricsDto, err := s.provider.FindLyrics(ctx, req)
		if err != nil {
			s.logger.Warnf("Error fetching lyrics for '%s': %v", item.SongTitle, err)
			continue
		}

		if lyricsDto != nil && lyricsDto.PlainLyrics != "" {
			text := lyricsDto.PlainLyrics
			if err := s.repo.UpdateStatusAndText(ctx, item.ID, StatusFound, &text); err != nil {
				s.logger.Errorf("Failed to update lyrics status to FOUND for song %s: %v", item.SongID, err)
			} else {
				s.logger.Infof("Lyrics found and saved for song: title='%s', artist='%s'", item.SongTitle, artist)
			}
		} else {
			// Выставляем FAILED, чтобы не повторять бесконечно запросы для трека без текста
			if err := s.repo.UpdateStatusAndText(ctx, item.ID, StatusFailed, nil); err != nil {
				s.logger.Errorf("Failed to update lyrics status to FAILED for song %s: %v", item.SongID, err)
			} else {
				s.logger.Warnf("Lyrics not found for song: title='%s', artist='%s' (status set to FAILED)", item.SongTitle, artist)
			}
		}
	}

	return nil
}
