package com.localsound.backend.download;

import jakarta.validation.constraints.NotBlank;

public record DownloadRequest(
	@NotBlank(message = "Url is required")
	String url
) {}
