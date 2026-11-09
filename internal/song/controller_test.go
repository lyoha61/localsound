package song_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"localsound-backend/internal/song"
	"localsound-backend/internal/ytdlp"
)

type mockService struct {
	songs         []song.SongDto
	searchResults []ytdlp.SearchResultItem
	err           error
}

func (m *mockService) GetSongs(ctx context.Context) ([]song.SongDto, error) {
	return m.songs, m.err
}

func (m *mockService) SearchSongs(ctx context.Context, query string) ([]ytdlp.SearchResultItem, error) {
	return m.searchResults, m.err
}

func TestController_GetSongs(t *testing.T) {
	svc := &mockService{
		songs: []song.SongDto{
			{Title: "Track 1"},
		},
	}

	testLogger := zap.NewNop().Sugar()

	ctrl := song.NewController(svc, testLogger)

	req := httptest.NewRequest(http.MethodGet, "/songs", nil)
	rec := httptest.NewRecorder()

	ctrl.GetSongs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}
