package com.localsound.backend.lyrics;

import java.util.List;
import java.util.Objects;
import java.util.Optional;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.RestClient;

import com.localsound.backend.lyrics.dto.LyricsDto;
import com.localsound.backend.lyrics.dto.LyricsRequest;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@Slf4j
public class LrcLibLyricsService implements LyricsProvider {
	private static final String BASE_URL = "https://lrclib.net/api";
	private final RestClient restClient;

	private Optional<LyricsDto> findByGet(LyricsRequest request) {
		var response = restClient.get()
			.uri(uriBuilder -> uriBuilder
				// .path(BASE_URL + "/get")
				.path("/get")
				.queryParam("track_name", request.trackName())
				.queryParam("artist_name", request.artistName())
				.queryParam("album_name", request.albumName())
				.queryParam("duration", request.duration())
				.build())
			// .header(
			// 	HttpHeaders.USER_AGENT, 
			// 	"LocalSound/1.0 (https://github.com/lyoha61/localsound)"
			// )
			.retrieve()
			.body(LyricsDto.class);
		
		return Optional.ofNullable(response);
	}

	private Optional<LyricsDto> findBySearch(LyricsRequest request) {
		log.info(
			"Searching LRCLIB via /search: track='{}', artist='{}'",
			request.trackName(),
			request.artistName()
    );

		var songParams = Stream.of(
			request.trackName(),
			request.artistName(),
			request.albumName()
		).filter(Objects::nonNull)
		.filter(s -> !s.isBlank())
		.collect(Collectors.joining(" "));

		var response = restClient.get()
			.uri(uriBuilder -> uriBuilder
				// .path(BASE_URL + "/search")
				.path("/search")
				.queryParam("q", songParams)
				.queryParam("track_name", request.trackName())
				.queryParam("artist_name", request.artistName())
				.queryParam("album_name", request.albumName())
				.build())
			// .header(
			// 	HttpHeaders.USER_AGENT, 
			// 	"LocalSound/1.0 (https://github.com/lyoha61/localsound)"
			// )
			.retrieve()
			.body(new ParameterizedTypeReference<List<LyricsDto>>() {});
		
		log.info("LRCLIB /search returned {} results", response.size());

		return response.stream().findFirst();
	}

	@Override
	public Optional<LyricsDto> findLyrics(LyricsRequest request) {
		try {
			return findByGet(request);
		} catch (HttpClientErrorException.NotFound e) {

			log.warn(
				"Lyrics not found via GET: track='{}', artist='{}', status={}",
				request.trackName(),
				request.artistName(),
				e.getStatusCode()
			);

			try {
				return findBySearch(request);
			} catch (HttpClientErrorException.NotFound e2) {
				log.error(
					"LRCLIB search failed: track='{}', artist='{}', status={}",
					request.trackName(),
					request.artistName(),
					e2.getStatusCode()
				);

				return Optional.empty();
			} 
		} catch (HttpClientErrorException.TooManyRequests e) {
			log.error(
				"LRCLIB rate limit exceeded: track='{}', artist='{}', status={}",
				request.trackName(),
				request.artistName(),
				e.getStatusCode()
			);
			
			return Optional.empty();
		} catch (HttpClientErrorException e) {
			log.error(
				"LRCLIB request failed: track='{}', artist='{}', status={}",
				request.trackName(),
				request.artistName(),
				e.getStatusCode(),
				e
			);

			return Optional.empty();
		}
	}
}
