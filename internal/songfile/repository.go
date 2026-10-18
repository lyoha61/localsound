package songfile

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

func (r *Repository) Create(ctx context.Context, f *SongFileEntity) error {
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now

	query := `
		INSERT INTO song_files (id, song_id, file_path, format, file_size, created_at, updated_at)
		VALUES (uuidv7(), :song_id, :file_path, :format, :file_size, :created_at, :updated_at)
		RETURNING id
	`

	rows, err := r.db.NamedQueryContext(ctx, query, f)
	if err != nil {
		return fmt.Errorf("failed to insert song file: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan generated song_file id: %w", err)
		}
		f.ID = id
	}

	return nil
}
