package com.localsound.backend.song;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

public interface SongRepository extends JpaRepository<SongEntity, UUID> {}
