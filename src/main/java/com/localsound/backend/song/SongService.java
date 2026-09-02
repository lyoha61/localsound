package com.localsound.backend.song;

import java.nio.file.Path;

import org.springframework.stereotype.Service;

import com.localsound.backend.lyrics.LyricsService;
import com.localsound.backend.song.dto.SongMetadata;
import com.localsound.backend.songfile.SongFileEntity;
import com.localsound.backend.songfile.SongFileRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class SongService {
	private final SongRepository songRepository;
	private final SongFileRepository songFileRepository;
	private final SongMetadataService songMetadataService;
	private final LyricsService lyricsService;

	@Transactional
	public void create(Path audioPath) {
		SongMetadata metadata = songMetadataService.getData(audioPath);

		var format = SongFormat.valueOf(metadata.format().toUpperCase());

		var songEntity = new SongEntity(
			metadata.title(),
			metadata.artist(),
			metadata.album(),
			metadata.duration()
		);

		var songFileEntity = new SongFileEntity(
			songEntity,
			audioPath.toString(),
			format,
			metadata.fileSize()
		);

		songRepository.save(songEntity);

		songFileRepository.save(songFileEntity);

		lyricsService.create(songEntity);

		log.info("Audio file {} saved in DB", audioPath.getFileName());
	}
}
