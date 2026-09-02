package com.localsound.backend.lyrics;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import lombok.RequiredArgsConstructor;

@Component
@RequiredArgsConstructor
public class LyricsWorker {
	
	private final LyricsService lyricsService;

	@Scheduled(fixedDelay = 10_000)
	public void process() {
		lyricsService.processPending();
	}
}
