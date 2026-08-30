package com.localsound.backend.lyrics;

import org.springframework.stereotype.Service;

import com.localsound.backend.song.SongEntity;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class LyricsService {
	private final LyricsRepository lyricsRepository;

	public void create(SongEntity songEntity) {
		var lyricsEntity = new LyricsEntity(songEntity);

		lyricsRepository.save(lyricsEntity);

		log.info("Lyrics saved for song id={}", songEntity.getId());
	}
}
