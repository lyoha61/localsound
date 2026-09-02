package com.localsound.backend.lyrics;

import java.util.Optional;

import com.localsound.backend.lyrics.dto.LyricsDto;
import com.localsound.backend.lyrics.dto.LyricsRequest;

public interface LyricsProvider {
	
	Optional<LyricsDto> findLyrics(LyricsRequest request);
}
