package com.localsound.backend.download;

import java.io.IOException;
import java.nio.file.Path;

import org.springframework.stereotype.Service;

import com.localsound.backend.song.SongService;
import com.localsound.backend.storage.SongStorageService;
import com.localsound.backend.ytdlp.YtdlpService;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class DownloadService {
	private final YtdlpService ytdlpService;
	private final SongStorageService songStorageService;
	private final SongService songService;
	
	public Path download(
		String url, 
		DownloadOptions options
	) throws IOException, InterruptedException {
		Path outputDirectory = songStorageService.getSongsDirectory();

		Path audioPath = ytdlpService.download(url, options, outputDirectory);

		songService.create(audioPath);

		log.info("URL processed: {}", url);

		return audioPath;
	}
}
