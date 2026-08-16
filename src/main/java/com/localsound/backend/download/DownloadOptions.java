package com.localsound.backend.download;

public record DownloadOptions(
	boolean extractAudio,
	String audioFormat
) {}
