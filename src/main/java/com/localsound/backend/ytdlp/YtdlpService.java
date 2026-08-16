package com.localsound.backend.ytdlp;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.localsound.backend.download.DownloadOptions;

@Service
public class YtdlpService {
	
	private static final Logger log = 
		LoggerFactory.getLogger(YtdlpService.class);

	private List<String> buildCommand(
		String url,
		DownloadOptions options,
		Path outputDirectory
	) {
		List<String> command = new ArrayList<>();

		command.add("yt-dlp");

		command.add("--cookies-from-browser");
		command.add("chrome");

		command.add("--no-playlist");

		if (options.extractAudio()) {
			command.add("-x");

			command.add("--audio-format");
			command.add(options.audioFormat());
		}

		command.add("--embed-metadata");

		command.add("--embed-thumbnail");
		command.add("--convert-thumbnails");
		command.add("jpg");

		command.add("-o");
		command.add(
			outputDirectory
				.resolve("%(title)s.%(ext)s")
				.toString()
		);

		command.add(url);

		return command;
	}

	public Path download(
		String url,
		DownloadOptions options,
		Path outputDirectory
	) throws IOException, InterruptedException {
		List<String> command = buildCommand(url, options, outputDirectory);

		var processBuiler = new ProcessBuilder(command);

		processBuiler.redirectErrorStream(true);

		var process  = processBuiler.start();

		try (var reader = process.inputReader()) {
			String line;

			while ((line = reader.readLine()) != null) {
				log.info("yt-dlp: {}", line);
			}
		}

		int exitcode = process.waitFor();

		if (exitcode != 0) {
			log.error("yt-dlp failed with exit code: {}", exitcode);
		}

		try (var files = Files.list(outputDirectory)) {
			return files
				.filter(Files::isRegularFile)
				.filter(path -> path.toString().endsWith("mp3"))
				.findFirst()
				.orElseThrow(() -> new IOException("Downloaded MP3 not found"));
		}
	}
}
