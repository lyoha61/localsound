package com.localsound.backend.storage;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.springframework.stereotype.Service;

import jakarta.annotation.PostConstruct;

@Service
public class TemporaryStorageService {

	private final Path downloadsDirectory = 
		Path.of("tmp", "downloads");
	
	@PostConstruct
	void init() throws IOException {
		Files.createDirectories(downloadsDirectory);
	}

	public Path createDownloadDirectory() throws IOException {
		return Files.createTempDirectory(
			downloadsDirectory, 
			"download-"
		);
	}
}
