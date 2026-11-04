package song

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

func (r *Repository) Create(ctx context.Context, s *SongEntity) error {
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now

	query := `
		INSERT INTO songs (id, title, artist, album, duration, created_at, updated_at)
		VALUES (uuidv7(), :title, :artist, :album, :duration, :created_at, :updated_at)
		RETURNING id
	`

	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("failed to insert song: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan generated song id: %w", err)
		}
		s.ID = id
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*SongEntity, error) {
	var s SongEntity
	query := `SELECT id, title, artist, album, duration, created_at, updated_at FROM songs WHERE id = $1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]SongEntity, error) {
	s := []SongEntity{}
	query := `SELECT id, title, artist, album, duration, created_at, updated_at FROM songs ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &s, query); err != nil {
		return nil, err
	}

	return s, nil
}
