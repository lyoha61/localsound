package com.localsound.backend.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpHeaders;
import org.springframework.web.client.RestClient;

@Configuration
public class RestClientConfig {

	    @Bean
    public RestClient restClient() {
        return RestClient.builder()
            .baseUrl("https://lrclib.net/api")
						.defaultHeader(
							HttpHeaders.USER_AGENT, 
							"LocalSound/1.0 (https://github.com/lyoha61/localsound)"
						)
            .build();
    }
}