package com.localsound.backend.songfile;

import java.nio.file.Path;
import java.time.Instant;
import java.util.UUID;

import com.localsound.backend.song.SongEntity;
import com.localsound.backend.song.SongFormat;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToOne;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.Setter;

@Entity
@Table(name = "song_files")
@Getter
@Setter
public class SongFileEntity {
	
	@Id
	@GeneratedValue
	UUID id;

	@OneToOne
	@JoinColumn(name = "song_id", nullable = false, unique = true)
	SongEntity song;

	@Column(nullable = false)
	private String filePath;

	@Enumerated(EnumType.STRING)
	@Column(nullable = false)
	private SongFormat format;

	@Column(nullable = false)
	private Long fileSize;

	@Column(nullable = false)
	private Instant createdAt;

	@Column(nullable = false)
	private Instant updatedAt;

	@PrePersist
	private void onCreate() {
		var now = Instant.now();

		this.createdAt = now;
		this.updatedAt = now;
	}

	@PreUpdate
	private void onUpdate() {
		this.updatedAt = Instant.now();
	}

	protected SongFileEntity() {}

	public SongFileEntity(
		SongEntity song,
		String filePath,
		SongFormat format,
		Long fileSize
	) {
		this.song = song;
		this.filePath = filePath;
		this.format = format;
		this.fileSize = fileSize; 
	}
}


	// song_id UUID NOT NULL UNIQUE REFERENCES songs(id) ON DELETE CASCADE,

	// file_path TEXT NOT NULL,
	// format VARCHAR(20) NOT NULL,
	// file_size BIGINT NOT NULL