package ytdlp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Platform string

const (
	PlatformSoundCloud Platform = "soundcloud"
	PlatformYouTube    Platform = "youtube"
)

type DownloadOptions struct {
	ExtractAudio bool
	AudioFormat  string
}

type Service struct {
	cookiesFromBrowser bool
	cookiesPath        string
	logger             *zap.SugaredLogger
}

type SearchResultItem struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Source    Platform `json:"source"`
	Uploader  string   `json:"uploader"`
	Duration  float64  `json:"duration"`
	URL       string   `json:"url"`
	Thumbnail string   `json:"thumbnail"`
}

type rawThumbnail struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type rawEntry struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Uploader   string         `json:"uploader"`
	Duration   float64        `json:"duration"`
	WebpageURL string         `json:"webpage_url"`
	Thumbnails []rawThumbnail `json:"thumbnails"`
}

type rawSearchResponse struct {
	Entries []rawEntry `json:"entries"`
}

func NewService(cookiesFromBrowser bool, cookiesPath string, logger *zap.SugaredLogger) *Service {
	return &Service{
		cookiesFromBrowser: cookiesFromBrowser,
		cookiesPath:        cookiesPath,
		logger:             logger,
	}
}

func (s *Service) addCookiesArguments(command []string) []string {
	if s.cookiesFromBrowser {
		return append(command, "--cookies-from-browser", "chrome")
	}

	if s.cookiesPath != "" {
		if _, err := os.Stat(s.cookiesPath); err == nil {
			return append(command, "--cookies", s.cookiesPath)
		}
	}

	return command
}

func (s *Service) buildCommand(url string, options DownloadOptions, outputDir string) []string {
	command := []string{
		"--js-runtimes", "node",
	}

	command = s.addCookiesArguments(command)
	command = append(command, "--no-playlist")

	if options.ExtractAudio {
		command = append(command, "-x")
		format := options.AudioFormat
		if format == "" {
			format = "mp3"
		}
		command = append(command, "--audio-format", format)
	}

	command = append(command,
		"--embed-metadata",
		"--embed-thumbnail",
		"--convert-thumbnails", "jpg",
		"-o", filepath.Join(outputDir, "%(title)s.%(ext)s"),
		"--print", "after_move:filepath",
		url,
	)

	return command
}

func (s *Service) exec(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		s.logger.Errorf("yt-dlp execution failed: %v, stderr: %s", err, stderrBuf.String())
		return nil, fmt.Errorf("yt-dlp failed: %w (stderr: %s)", err, strings.TrimSpace(stderrBuf.String()))
	}

	return stdoutBuf.Bytes(), nil
}

func (s *Service) Download(ctx context.Context, url string, options DownloadOptions, outputDir string) (string, error) {
	execCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	args := s.buildCommand(url, options, outputDir)
	cmd := exec.CommandContext(execCtx, "yt-dlp", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	s.logger.Infof("Executing yt-dlp for URL: %s (format: %s)", url, options.AudioFormat)

	err := cmd.Run()
	if err != nil {
		s.logger.Errorf("yt-dlp execution failed: %v, stderr: %s", err, stderrBuf.String())
		return "", fmt.Errorf("yt-dlp failed: %w (stderr: %s)", err, strings.TrimSpace(stderrBuf.String()))
	}

	stdoutLines := strings.Split(strings.TrimSpace(stdoutBuf.String()), "\n")
	var audioPath string
	for i := len(stdoutLines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(stdoutLines[i])
		if line != "" {
			if _, statErr := os.Stat(line); statErr == nil {
				audioPath = line
				break
			}
		}
	}

	if audioPath == "" {
		return "", fmt.Errorf("downloaded audio file not found in yt-dlp output: %s", stdoutBuf.String())
	}

	s.logger.Infof("Successfully downloaded audio: %s", audioPath)
	return audioPath, nil
}

func (s *Service) SearchAll(ctx context.Context, query string, limitPerSource int) ([]SearchResultItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type searchJob struct {
		source Platform
		target string
	}

	jobs := []searchJob{
		{
			source: PlatformSoundCloud,
			target: fmt.Sprintf("scsearch%d:%s", limitPerSource, query),
		},
		{
			source: PlatformYouTube,
			target: fmt.Sprintf("ytsearch%d:%s", limitPerSource, query),
		},
	}

	resultsChan := make(chan []SearchResultItem, len(jobs))
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go func(j searchJob) {
			defer wg.Done()

			items, err := s.searchTarget(ctx, j.source, j.target)
			if err != nil {
				s.logger.Warnf("Search failed for %s: %v", j.source, err)
				return
			}
			resultsChan <- items
		}(job)
	}

	wg.Wait()
	close(resultsChan)

	var allResults []SearchResultItem
	for items := range resultsChan {
		allResults = append(allResults, items...)
	}
	return allResults, nil
}

func (s *Service) normalizeURL(platform Platform, webpageURL, id string) string {
	if webpageURL != "" && strings.HasPrefix(webpageURL, "http") {
		return webpageURL
	}

	switch platform {
	case PlatformYouTube:
		return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
	case PlatformSoundCloud:
		return webpageURL
	default:
		return webpageURL
	}
}

func (s *Service) searchTarget(ctx context.Context, platform Platform, target string) ([]SearchResultItem, error) {
	args := []string{
		// stdout with format JSON
		"--dump-single-json",
		// Return just list songs
		"--flat-playlist",
		"--no-update",
		target,
	}

	stdoutBytes, err := s.exec(ctx, args)
	if err != nil {
		return nil, err
	}

	var rawResp rawSearchResponse
	if err := json.Unmarshal(stdoutBytes, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to decode yt-dlp json: %w", err)
	}

	results := make([]SearchResultItem, 0, len(rawResp.Entries))

	for _, entry := range rawResp.Entries {

		bestThumb := ""
		for _, th := range entry.Thumbnails {
			if th.ID == "t500x500" {
				bestThumb = th.URL
				break
			}
			bestThumb = th.URL
		}

		results = append(results, SearchResultItem{
			ID:        entry.ID,
			Source:    platform,
			Title:     entry.Title,
			Uploader:  entry.Uploader,
			Duration:  entry.Duration,
			URL:       s.normalizeURL(platform, entry.WebpageURL, entry.ID),
			Thumbnail: bestThumb,
		})
	}
	return results, nil
}
