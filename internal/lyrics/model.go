package lyrics

import (
	"time"

	"github.com/google/uuid"
)

type LyricsStatus string

const (
	StatusPending    LyricsStatus = "PENDING"
	StatusProcessing LyricsStatus = "PROCESSING"
	StatusFound      LyricsStatus = "FOUND"
	StatusFailed     LyricsStatus = "FAILED"
)

type LyricsEntity struct {
	ID           uuid.UUID    `db:"id" json:"id"`
	SongID       uuid.UUID    `db:"song_id" json:"songId"`
	Text         *string      `db:"text" json:"text"`
	LyricsStatus LyricsStatus `db:"lyrics_status" json:"lyricsStatus"`
	SyncStatus   string       `db:"sync_status" json:"syncStatus"`
	CreatedAt    time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updatedAt"`
}

type LyricsWithSong struct {
	LyricsEntity
	SongTitle    string  `db:"song_title"`
	SongArtist   *string `db:"song_artist"`
	SongAlbum    *string `db:"song_album"`
	SongDuration int     `db:"song_duration"`
}

type LyricsRequest struct {
	TrackName  string
	ArtistName string
	AlbumName  string
	Duration   int
}

type LyricsDto struct {
	ID           int64   `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}
