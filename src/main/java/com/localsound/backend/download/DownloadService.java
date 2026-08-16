package com.localsound.backend.download;

import java.io.IOException;
import java.nio.file.Path;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.stereotype.Service;

import com.localsound.backend.storage.TemporaryStorageService;
import com.localsound.backend.ytdlp.YtdlpService;

@Service
public class DownloadService {
	private final YtdlpService ytdlpService;
	private final TemporaryStorageService temporaryStorageService;

	public DownloadService(
		YtdlpService ytdlpService,
		TemporaryStorageService temporaryStorageService
	) {
		this.ytdlpService = ytdlpService;
		this.temporaryStorageService = temporaryStorageService;
	}
	
	private static final Logger log = 
		LoggerFactory.getLogger(DownloadService.class);

	public Path download(
		String url, 
		DownloadOptions options
	) throws IOException, InterruptedException {
		Path outputDirectory = temporaryStorageService.createDownloadDirectory();

		Path audioPath = ytdlpService.download(url, options, outputDirectory);

		log.info("URL processed: {}", url);

		return audioPath;
	}
}
