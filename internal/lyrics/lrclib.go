package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Provider interface {
	FindLyrics(ctx context.Context, req LyricsRequest) (*LyricsDto, error)
	Search(ctx context.Context, query string) ([]LyricsDto, error)
}

type LrcLibClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.SugaredLogger
}

type userAgentTransport struct {
	rt        http.RoundTripper
	userAgent string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.rt.RoundTrip(req)
}

func NewLrcLibClient(baseURL string, logger *zap.SugaredLogger) *LrcLibClient {
	return &LrcLibClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &userAgentTransport{
				rt:        http.DefaultTransport,
				userAgent: "LocalSound/1.0 (https://github.com/lyoha61/localsound)",
			},
		},
		logger: logger,
	}
}

func (c *LrcLibClient) FindLyrics(ctx context.Context, req LyricsRequest) (*LyricsDto, error) {
	lyrics, err := c.findByGet(ctx, req)
	if err == nil && lyrics != nil {
		return lyrics, nil
	}

	c.logger.Infof("Lyrics not found via GET, trying search for track='%s', artist='%s'", req.TrackName, req.ArtistName)
	return c.findBySearch(ctx, req)
}

func (c *LrcLibClient) findByGet(ctx context.Context, req LyricsRequest) (*LyricsDto, error) {
	params := url.Values{}
	params.Set("track_name", req.TrackName)
	if req.ArtistName != "" {
		params.Set("artist_name", req.ArtistName)
	}
	if req.AlbumName != "" {
		params.Set("album_name", req.AlbumName)
	}
	if req.Duration > 0 {
		params.Set("duration", strconv.Itoa(req.Duration))
	}

	endpoint := fmt.Sprintf("%s/get?%s", c.baseURL, params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LRCLIB /get request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LRCLIB /get status code %d", resp.StatusCode)
	}

	var dto LyricsDto
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("failed to decode LRCLIB response: %w", err)
	}

	return &dto, nil
}

func (c *LrcLibClient) findBySearch(ctx context.Context, req LyricsRequest) (*LyricsDto, error) {
	var queryParts []string
	if req.TrackName != "" {
		queryParts = append(queryParts, req.TrackName)
	}
	if req.ArtistName != "" {
		queryParts = append(queryParts, req.ArtistName)
	}
	if req.AlbumName != "" {
		queryParts = append(queryParts, req.AlbumName)
	}

	params := url.Values{}
	params.Set("q", strings.Join(queryParts, " "))
	if req.TrackName != "" {
		params.Set("track_name", req.TrackName)
	}
	if req.ArtistName != "" {
		params.Set("artist_name", req.ArtistName)
	}
	if req.AlbumName != "" {
		params.Set("album_name", req.AlbumName)
	}

	endpoint := fmt.Sprintf("%s/search?%s", c.baseURL, params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LRCLIB /search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LRCLIB /search status code %d", resp.StatusCode)
	}

	var list []LyricsDto
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode LRCLIB search response: %w", err)
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &list[0], nil
}

func (c *LrcLibClient) Search(ctx context.Context, query string) ([]LyricsDto, error) {
	endpoint := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return nil, err
	}

	c.logger.Infof("Sending request to LRCLIB: query='%s'", query)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LRCLIB /search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LRCLIB /search status code %d", resp.StatusCode)
	}

	var list []LyricsDto
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode LRCLIB search response: %w", err)
	}

	c.logger.Infof("LRCLIB search finished: found %d results for query='%s'", len(list), query)

	return list, nil
}
