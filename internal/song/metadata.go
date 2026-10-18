package song

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.senan.xyz/taglib"
)

type MetadataService struct{}

func NewMetadataService() *MetadataService {
	return &MetadataService{}
}

func (s *MetadataService) GetData(audioPath string) (*SongMetadata, error) {
	fi, err := os.Stat(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio file info %s: %w", audioPath, err)
	}

	tags, err := taglib.ReadTags(audioPath)
	if err != nil {
		// Даже если теги повреждены, пытаемся продолжить
		tags = make(map[string][]string)
	}

	props, err := taglib.ReadProperties(audioPath)
	if err != nil {
		props = taglib.Properties{}
	}

	firstTag := func(key string) string {
		vals, ok := tags[key]
		if ok && len(vals) > 0 {
			return strings.TrimSpace(vals[0])
		}
		return ""
	}

	title := firstTag(taglib.Title)
	if title == "" {
		// Если в тегах нет заголовка, используем имя файла без расширения
		fileName := filepath.Base(audioPath)
		ext := filepath.Ext(fileName)
		title = strings.TrimSuffix(fileName, ext)
	}

	artist := firstTag(taglib.Artist)
	album := firstTag(taglib.Album)

	duration := int(props.Length.Seconds())
	format := props.Format
	if format == "" {
		format = strings.TrimPrefix(filepath.Ext(audioPath), ".")
	}

	return &SongMetadata{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
		Format:   format,
		FileSize: fi.Size(),
	}, nil
}
