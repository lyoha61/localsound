package com.localsound.backend.song;

import java.time.Instant;
import java.util.UUID;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.Id;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.Setter;

@Entity
@Table(name = "songs")
@Getter
@Setter
public class SongEntity {

	@Id
	@GeneratedValue
	private UUID id;

	@Column(nullable = false)
	private String title;

	@Column(nullable = true)
	private String artist;
	
	@Column(nullable = true)
	private String album;

	@Column(nullable = false)
	private Integer duration;

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

	protected SongEntity(){}

	public SongEntity(
		String title,
		String artist,
		String album,
		Integer duration,
		String filePath,
		SongFormat format,
		Long fileSize
	) {
		this.title = title;
		this.artist = artist;
		this.album = album;
		this.duration = duration;
		this.filePath = filePath;
		this.format = format;
		this.fileSize = fileSize;
	}
}