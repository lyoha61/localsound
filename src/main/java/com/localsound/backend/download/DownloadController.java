package com.localsound.backend.download;

import java.io.IOException;
import java.nio.file.Path;

import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import jakarta.validation.Valid;

@RestController
@RequestMapping("/download")
public class DownloadController {
	
	private final DownloadService downloadService;

	public DownloadController(DownloadService downloadService) {
		this.downloadService = downloadService;
	}

	@PostMapping()
	public ResponseEntity<Resource> downloadFromUrl(
		@Valid @RequestBody DownloadRequest body
	) throws IOException, InterruptedException {

		DownloadOptions options = new DownloadOptions(
			true,
			"mp3"
		);

		Path audioPath = downloadService.download(body.url(), options);

		Resource resource = new FileSystemResource(audioPath);

		return ResponseEntity.ok()
			.contentType(MediaType.parseMediaType("audio/mpeg"))
			.header(
				HttpHeaders.CONTENT_DISPOSITION, 
				"attachment; filename=\"" + audioPath.getFileName() + "\""
			)
			.body(resource);
	}
}
