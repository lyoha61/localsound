package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"localsound-backend/internal/songfile"
	"localsound-backend/internal/ytdlp"

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

type DownloadRequest struct {
	URL    string              `json:"url"`
	Format songfile.SongFormat `json:"format"`
}

func (c *Controller) DownloadFromURL(w http.ResponseWriter, r *http.Request) {
	var body DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}

	if !body.Format.IsValid() {
		http.Error(w, `{"error": "Invalid format"}`, http.StatusBadRequest)
		return
	}

	options := ytdlp.DownloadOptions{
		ExtractAudio: true,
		AudioFormat:  strings.ToLower(string(body.Format)),
	}

	audioPath, err := c.service.Download(r.Context(), body.URL, options)
	if err != nil {
		c.logger.Errorf("Download failed for URL %s: %v", body.URL, err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	file, err := os.Open(audioPath)
	if err != nil {
		c.logger.Errorf("Failed to open downloaded audio file %s: %v", audioPath, err)
		http.Error(w, `{"error": "Failed to read downloaded audio file"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileName := filepath.Base(audioPath)
	encodedFileName := url.PathEscape(fileName)

	if stat, err := file.Stat(); err == nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	w.Header().Set("Content-Type", body.Format.ContentType())
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, encodedFileName))

	c.logger.Infof("Starting stream: %s", audioPath)

	if _, err := io.Copy(w, file); err != nil {
		c.logger.Errorf("Error during file streaming %s: %v", audioPath, err)
		return
	}

	c.logger.Infof("Finished stream: %s", audioPath)
}
