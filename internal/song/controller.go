package song

import (
	"context"
	"encoding/json"
	"fmt"
	"localsound-backend/internal/httputil"
	"localsound-backend/internal/response"
	"localsound-backend/internal/songfile"
	"localsound-backend/internal/ytdlp"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SongService interface {
	GetSongs(ctx context.Context) ([]SongDto, error)
	SearchSongs(ctx context.Context, query string) ([]ytdlp.SearchResultItem, error)
}

type Controller struct {
	service SongService
	logger  *zap.SugaredLogger
}

type SongDto struct {
	Title     string    `db:"title" json:"title"`
	Artist    *string   `db:"artist" json:"artist"`
	Album     *string   `db:"album" json:"album"`
	Duration  int       `db:"duration" json:"duration"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type ImportTrackRequest struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type ImportBatchRequest struct {
	Format songfile.SongFormat  `json:"format"`
	Tracks []ImportTrackRequest `json:"tracks"`
}

func NewController(service SongService, logger *zap.SugaredLogger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c *Controller) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := c.service.GetSongs(r.Context())
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "failed to get songs")
		return
	}

	_ = response.JSON(w, http.StatusOK, songs)
}

func (c *Controller) ImportSongs(w http.ResponseWriter, r *http.Request) {
	var req ImportBatchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if len(req.Tracks) == 0 {
		http.Error(w, "Tracks list cannot be empty", http.StatusBadRequest)
		return
	}

	if !req.Format.IsValid() {
		http.Error(w, "Invalid audio format", http.StatusBadRequest)
		return
	}

	for i, track := range req.Tracks {
		if strings.TrimSpace(track.Artist) == "" ||
			strings.TrimSpace(track.Title) == "" {
			http.Error(w, fmt.Sprintf("Track at index %d must have non-empty artist and title", i), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) SearchSong(w http.ResponseWriter, r *http.Request) {
	query, ok := httputil.GetRequiredQuery(w, r, "q")
	if !ok {
		return
	}

	results, err := c.service.SearchSongs(r.Context(), query)
	if err != nil {
		c.logger.Errorf("Failed to search songs: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "failed to search songs")
		return
	}

	_ = response.JSON(w, http.StatusOK, results)
}
