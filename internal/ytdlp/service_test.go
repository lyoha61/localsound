package ytdlp

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestService_Search_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	logger := zap.NewNop().Sugar()
	svc := NewService(false, "", logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := "Гуф Не все потеряно"
	limitPerSource := 3

	results, err := svc.SearchAll(ctx, query, limitPerSource)
	if err != nil {
		t.Fatalf("SearchAll failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least 1 track, got 0")
	}

	t.Logf("Found %d tracks across platforms:", len(results))
	for i, track := range results {
		t.Logf("[%d] [%s] %s — %s (%.0fs) | URL: %s",
			i+1,
			track.Source,
			track.Uploader,
			track.Title,
			track.Duration,
			track.URL,
		)
	}
}
