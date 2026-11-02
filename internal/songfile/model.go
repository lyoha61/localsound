package songfile

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type SongFormat string

const (
	FormatMP3 SongFormat = "MP3"
	FormatM4A SongFormat = "M4A"
)

func (sf SongFormat) IsValid() bool {
	switch strings.ToUpper(string(sf)) {
	case string(FormatMP3), string(FormatM4A):
		return true
	default:
		return false
	}
}

func (sf SongFormat) ContentType() string {
	switch strings.ToUpper(string(sf)) {
	case string(FormatM4A):
		return "audio/mp4"
	default:
		return "audio/mpeg"
	}
}

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
