package com.localsound.backend.lyrics;

import java.util.Optional;

import org.springframework.stereotype.Service;

import com.localsound.backend.lyrics.dto.LyricsDto;
import com.localsound.backend.lyrics.dto.LyricsRequest;
import com.localsound.backend.song.SongEntity;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class LyricsService {
	private final LyricsRepository lyricsRepository;
	private final LyricsProvider lyricsProvider;

	public void create(SongEntity songEntity) {
		var lyricsEntity = new LyricsEntity(songEntity);

		lyricsRepository.save(lyricsEntity);

		log.info("Lyrics saved for song id={}", songEntity.getId());
	}

	public void processPending() {
		var lyricsEntities = lyricsRepository.findByLyricsStatus(LyricsStatus.PENDING);

		log.info("Processing {} pending lyrics", lyricsEntities.size());

		for (var lyricsEntity: lyricsEntities) {
			try {
				Thread.sleep(5000);
			} catch (InterruptedException e) {
				Thread.currentThread().interrupt();
				return;
			}
			var song = lyricsEntity.getSong();

			log.info(
				"Searching lyrics for song: title='{}', artist='{}'",
				song.getTitle(),
				song.getArtist()
			);

			var request = new LyricsRequest(
				song.getTitle(), 
				song.getArtist(),
				song.getAlbum(),
				song.getDuration()
			);

			Optional<LyricsDto> lyrics = lyricsProvider.findLyrics(request);

			if (lyrics.isPresent()) {
				var lyricsDto = lyrics.get();

				lyricsEntity.setText(lyricsDto.plainLyrics());
				lyricsEntity.setLyricsStatus(LyricsStatus.FOUND);
				lyricsRepository.save(lyricsEntity);

				log.info(
					"Lyrics found and saved for song: title='{}', artist='{}'",
					song.getTitle(),
					song.getArtist()
				);
			} else {
				log.warn(
					"Lyrics not found for song: title='{}', artist='{}'",
					song.getTitle(),
					song.getArtist()
				);
			}
		}
	}
}
