package com.localsound.backend.storage;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.springframework.stereotype.Service;

import jakarta.annotation.PostConstruct;

@Service
public class SongStorageService {
	
	private final Path songsDirectory = 
		Path.of("storage", "songs");
	
	@PostConstruct
	void init() throws IOException {
		Files.createDirectories(songsDirectory);
	}

	public Path getSongsDirectory() {
		return songsDirectory;
	}
}
