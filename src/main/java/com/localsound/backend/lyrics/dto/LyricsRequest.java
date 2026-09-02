package com.localsound.backend.lyrics.dto;

public record LyricsRequest(
	String trackName,
	String artistName,
	String albumName,
	int duration
) {}
