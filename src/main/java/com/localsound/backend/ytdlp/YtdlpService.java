package com.localsound.backend.ytdlp;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import com.localsound.backend.download.DownloadOptions;

import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class YtdlpService {
	
	@Value("${youtube.cookies-from-browser}")
	private boolean cookiesFromBrowser;

	private void addCookiesArguments(List<String> command) {
		if (cookiesFromBrowser) {
			command.add("--cookies-from-browser");
			command.add("chrome");
		} else {
			command.add("--cookies");
			command.add("/app/cookies/youtube.txt");
		}
	}

	private List<String> buildCommand(
		String url,
		DownloadOptions options,
		Path outputDirectory
	) {
		List<String> command = new ArrayList<>();

		command.add("yt-dlp");

		command.add("--js-runtimes");
		command.add("node");

		addCookiesArguments(command);

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

		command.add("--print");
		command.add("after_move:filepath");

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

		Path audioPath = null;

		try (var reader = process.inputReader()) {
			String line;

			while ((line = reader.readLine()) != null) {
				log.info("yt-dlp: {}", line);

				if (line.endsWith(".mp3")) {
					audioPath = Path.of(line);
				}
			}
		}

		int exitcode = process.waitFor();

		if (exitcode != 0) {
			log.error("yt-dlp failed with exit code: {}", exitcode);
		}

		if (audioPath == null) {
			throw new IOException("Downloaded MP3 not found");
		}

		return audioPath;
	}
}
