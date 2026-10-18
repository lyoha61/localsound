package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"localsound-backend/internal/config"
	"localsound-backend/internal/download"
	"localsound-backend/internal/lyrics"
	"localsound-backend/internal/song"
	"localsound-backend/internal/songfile"
	"localsound-backend/internal/storage"
	"localsound-backend/internal/ytdlp"
	"localsound-backend/migrations"
)

func main() {
	// Инициализация Zap логгера
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	logger.Infof("Starting LocalSound Backend on port %s", cfg.Port)

	// Подключение к PostgreSQL
	db, err := sqlx.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения с базой (с повторными попытками при старте)
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()
	if err := db.PingContext(ctxPing); err != nil {
		logger.Warnf("Waiting for database connection (%v)...", err)
	}

	// Запуск миграций БД через Goose
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("Failed to set goose dialect: %v", err)
	}

	// Попытка применить миграции
	if err := goose.Up(db.DB, "."); err != nil {
		logger.Warnf("Goose migrations warning: %v", err)
	} else {
		logger.Info("Database migrations applied successfully")
	}

	// Инициализация сервисов и репозиториев
	storageSvc, err := storage.NewStorageService(cfg.StorageDir)
	if err != nil {
		logger.Fatalf("Failed to initialize storage: %v", err)
	}

	metadataSvc := song.NewMetadataService()
	songRepo := song.NewRepository(db)
	songFileRepo := songfile.NewRepository(db)
	lyricsRepo := lyrics.NewRepository(db)

	lrclibClient := lyrics.NewLrcLibClient(cfg.LrclibBaseURL, logger)
	lyricsSvc := lyrics.NewService(lyricsRepo, lrclibClient, logger)
	songSvc := song.NewService(songRepo, songFileRepo, metadataSvc, lyricsSvc, logger)

	ytdlpSvc := ytdlp.NewService(cfg.CookiesFromBrowser, cfg.CookiesPath, logger)
	downloadSvc := download.NewService(ytdlpSvc, storageSvc, songSvc, logger)
	downloadCtrl := download.NewController(downloadSvc, logger)

	// Запуск фонового воркера для текстов песен
	workerCtx, workerCancel := context.WithCancel(context.Background())
	lyricsWorker := lyrics.NewWorker(lyricsSvc, cfg.LrclibWorkerInterval, logger)
	go lyricsWorker.Start(workerCtx)

	// Настройка HTTP роутера Chi
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Маппинг эндпоинта
	r.Post("/download", downloadCtrl.DownloadFromURL)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Minute, // Большой таймаут для длительных операций скачивания
		WriteTimeout: 10 * time.Minute,
	}

	// Graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		logger.Infof("Server listening on http://localhost:%s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		logger.Infof("Shutdown signal received: %v. Initiating graceful shutdown...", sig)
		workerCancel()

		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(ctxShutdown); err != nil {
			logger.Errorf("Error during server shutdown: %v", err)
			_ = server.Close()
		}
	}

	logger.Info("LocalSound Backend gracefully stopped")
}
