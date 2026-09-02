package com.localsound.backend.songfile;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

public interface SongFileRepository extends JpaRepository<SongFileEntity, UUID> {
	
}
