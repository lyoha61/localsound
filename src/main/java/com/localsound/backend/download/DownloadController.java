package com.localsound.backend.download;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

import org.springframework.http.ContentDisposition;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody;

import jakarta.validation.Valid;

@RestController
@RequestMapping("/download")
public class DownloadController {
	
	private final DownloadService downloadService;

	public DownloadController(DownloadService downloadService) {
		this.downloadService = downloadService;
	}

	@PostMapping()
	public ResponseEntity<StreamingResponseBody> downloadFromUrl(
		@Valid @RequestBody DownloadRequest body
	) throws IOException, InterruptedException {

		DownloadOptions options = new DownloadOptions(
			true,
			"mp3"
		);

		Path audioPath = downloadService.download(body.url(), options);

		ContentDisposition contentDisposition = ContentDisposition
			.attachment()
			.filename(audioPath.getFileName().toString(), StandardCharsets.UTF_8)
			.build();

		StreamingResponseBody response = outputStream -> {
			try {
				Files.copy(audioPath, outputStream);
			} finally {
				Files.deleteIfExists(audioPath);
				Files.deleteIfExists(audioPath.getParent());
			}
		};

		return ResponseEntity.ok()
			.contentType(MediaType.parseMediaType("audio/mpeg"))
			.header(
				HttpHeaders.CONTENT_DISPOSITION, 
				contentDisposition.toString()
			)
			.body(response);
	}
}
