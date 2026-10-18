package lyrics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, l *LyricsEntity) error {
	now := time.Now().UTC()
	l.CreatedAt = now
	l.UpdatedAt = now

	query := `
		INSERT INTO lyrics (id, song_id, text, lyrics_status, sync_status, created_at, updated_at)
		VALUES (uuidv7(), :song_id, :text, :lyrics_status, :sync_status, :created_at, :updated_at)
		RETURNING id
	`

	rows, err := r.db.NamedQueryContext(ctx, query, l)
	if err != nil {
		return fmt.Errorf("failed to insert lyrics: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan generated lyrics id: %w", err)
		}
		l.ID = id
	}

	return nil
}

func (r *Repository) FindByStatus(ctx context.Context, status LyricsStatus) ([]LyricsWithSong, error) {
	query := `
		SELECT 
			l.id, l.song_id, l.text, l.lyrics_status, l.sync_status, l.created_at, l.updated_at,
			s.title AS song_title, s.artist AS song_artist, s.album AS song_album, s.duration AS song_duration
		FROM lyrics l
		JOIN songs s ON s.id = l.song_id
		WHERE l.lyrics_status = $1
		ORDER BY l.created_at ASC
	`
	var results []LyricsWithSong
	if err := r.db.SelectContext(ctx, &results, query, status); err != nil {
		return nil, fmt.Errorf("failed to query lyrics by status: %w", err)
	}
	return results, nil
}

func (r *Repository) UpdateStatusAndText(ctx context.Context, id uuid.UUID, status LyricsStatus, text *string) error {
	query := `
		UPDATE lyrics
		SET lyrics_status = $1, text = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, status, text, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update lyrics status and text: %w", err)
	}
	return nil
}
