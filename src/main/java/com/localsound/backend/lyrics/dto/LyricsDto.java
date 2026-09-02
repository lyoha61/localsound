package com.localsound.backend.lyrics.dto;

public record LyricsDto(
	String trackName,
	String artistName,
	String albumName,
	int duration,
	String plainLyrics,
	String syncedLyrics
) {}
