package song

import (
	"time"

	"github.com/google/uuid"
)

type SongEntity struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Artist    *string   `db:"artist" json:"artist"`
	Album     *string   `db:"album" json:"album"`
	Duration  int       `db:"duration" json:"duration"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type SongMetadata struct {
	Title    string
	Artist   string
	Album    string
	Duration int
	Format   string
	FileSize int64
}
