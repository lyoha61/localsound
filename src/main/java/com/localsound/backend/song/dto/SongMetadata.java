package com.localsound.backend.song.dto;

public record SongMetadata(
	String title,
	String artist,
	String album,
	int duration,
	String format,
	long fileSize
) {}
