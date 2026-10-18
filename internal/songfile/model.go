package songfile

import (
	"time"

	"github.com/google/uuid"
)

type SongFormat string

const (
	FormatMP3 SongFormat = "MP3"
	FormatM4A SongFormat = "M4A"
)

func ParseSongFormat(format string) SongFormat {
	switch format {
	case "M4A", "m4a", "mp4", "aac":
		return FormatM4A
	default:
		return FormatMP3
	}
}

type SongFileEntity struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	SongID    uuid.UUID  `db:"song_id" json:"songId"`
	FilePath  string     `db:"file_path" json:"filePath"`
	Format    SongFormat `db:"format" json:"format"`
	FileSize  int64      `db:"file_size" json:"fileSize"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
}
