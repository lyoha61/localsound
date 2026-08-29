package com.localsound.backend.song;

import java.nio.file.Path;

import org.springframework.stereotype.Service;

import com.localsound.backend.song.dto.SongMetadata;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class SongService {
	private final SongMetadataService songMetadataService;
	private final SongRepository songRepository;

	public void create(Path audioPath) {
		SongMetadata metadata = songMetadataService.getData(audioPath);

		var format = SongFormat.valueOf(metadata.format().toUpperCase());

		var songEntity = new SongEntity(
			metadata.title(),
			metadata.artist(),
			metadata.album(),
			metadata.duration(),
			audioPath.toString(),
			format,
			metadata.fileSize()
		);

		songRepository.save(songEntity);
		log.info("Audio file {} saved in DB", audioPath.getFileName());
	}
}
