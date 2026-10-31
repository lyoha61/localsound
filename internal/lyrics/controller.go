package lyrics

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type Controller struct {
	service *Service
	logger  *zap.SugaredLogger
}

func NewController(service *Service, logger *zap.SugaredLogger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c *Controller) SearchLyricsLRCLIB(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if strings.TrimSpace(query) == "" {
		http.Error(w, `{"error": "Query parameter 'q' is required"}`, http.StatusBadRequest)
		return
	}

	results, err := c.service.Search(r.Context(), query)
	if err != nil {
		c.logger.Errorf("Failed to search lyrics: %v", err)
		http.Error(w, `{"error": "Failed to search lyrics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
