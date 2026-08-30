package com.localsound.backend.lyrics;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

public interface LyricsRepository extends JpaRepository<LyricsEntity, UUID> {}
