package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
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

func (s *Service) Download(ctx context.Context, url string, options DownloadOptions, outputDir string) (string, error) {
	// Добавляем защитный таймаут 5 минут на выполнение yt-dlp, если в контексте нет меньшего таймаута
	execCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	args := s.buildCommand(url, options, outputDir)
	cmd := exec.CommandContext(execCtx, "yt-dlp", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	s.logger.Infof("Executing yt-dlp for URL: %s", url)

	err := cmd.Run()
	if err != nil {
		s.logger.Errorf("yt-dlp execution failed: %v, stderr: %s", err, stderrBuf.String())
		return "", fmt.Errorf("yt-dlp failed: %w (stderr: %s)", err, strings.TrimSpace(stderrBuf.String()))
	}

	// Извлекаем путь к файлу из stdout (благодаря опции --print after_move:filepath)
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
