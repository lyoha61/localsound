# LocalSound Backend (Go)

Backend для музыкального приложения LocalSound.

## Документация
* [Flow скачивания песни по URL](docs/song-download-flow.md)
* [Flow поиска текста песни](docs/lyrics-fetching-flow.md)

## Стек технологий
* **Язык:** Go 1.26
* **Роутер:** `github.com/go-chi/chi/v5`
* **Логирование:** `go.uber.org/zap`
* **База данных:** PostgreSQL (`github.com/jackc/pgx/v5`, `github.com/jmoiron/sqlx`)
* **Миграции:** `github.com/pressly/goose/v3`
* **Метаданные аудио:** `go.senan.xyz/taglib` (TagLib WASM/wazero, без CGO)
* **Загрузка и обработка аудио:** `yt-dlp` и `ffmpeg`
* **Контейнеризация:** Docker, Docker Compose

## Возможности
* Загрузка аудио по URL с поддержкой cookies
* Обработка, чтение и запись аудио-метаданных (ID3/Tags)
* Фоновый поиск и синхронизация текстов песен через LRCLIB
* Потоковая отдача аудиофайлов (streaming)

## Запуск в разработке (Dev)

### Через Docker Compose:
```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

### Локально (БД в Docker, приложение на хосте):
```bash
docker compose up -d postgres
go run ./cmd/api
```