package song

import (
	"context"
	"localsound-backend/internal/response"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type SongService interface {
	GetSongs(ctx context.Context) ([]SongDto, error)
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
