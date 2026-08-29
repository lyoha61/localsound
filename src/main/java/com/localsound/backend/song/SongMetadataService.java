package com.localsound.backend.song;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.jaudiotagger.audio.AudioFile;
import org.jaudiotagger.audio.AudioFileIO;
import org.jaudiotagger.audio.AudioHeader;
import org.jaudiotagger.audio.exceptions.CannotReadException;
import org.jaudiotagger.audio.exceptions.InvalidAudioFrameException;
import org.jaudiotagger.audio.exceptions.ReadOnlyFileException;
import org.jaudiotagger.tag.FieldKey;
import org.jaudiotagger.tag.Tag;
import org.jaudiotagger.tag.TagException;
import org.springframework.stereotype.Service;

import com.localsound.backend.song.dto.SongMetadata;

@Service
public class SongMetadataService {
	public SongMetadata getData(Path audioPath) {
		try {
			AudioFile audioFile = AudioFileIO.read(audioPath.toFile());
			AudioHeader audioHeader = audioFile.getAudioHeader();
			Tag tag = audioFile.getTag();
	
			String title = tag.getFirst(FieldKey.TITLE);
			String artist = tag.getFirst(FieldKey.ARTIST);
			String album = tag.getFirst(FieldKey.ALBUM);
			
			int duration = audioHeader.getTrackLength();
			String format = audioHeader.getFormat();

			long fileSize = Files.size(audioPath);

			return new SongMetadata(
				title,
				artist,
				album,
				duration,
				format,
				fileSize
			);
		} catch (CannotReadException |
						IOException |
						TagException |
						ReadOnlyFileException |
						InvalidAudioFrameException e) {
					
				throw new RuntimeException("Failed to read metadata " + audioPath, e);
		 }
	}
}
