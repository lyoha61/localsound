package com.localsound.backend.lyrics;

import java.time.Instant;
import java.util.UUID;

import com.localsound.backend.song.SongEntity;

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
@Table(name = "lyrics")
@Getter
@Setter
public class LyricsEntity {
	@Id
	@GeneratedValue
	private UUID id;

	@OneToOne
	@JoinColumn(name = "song_id", nullable = false, unique = true)
	private SongEntity song;

	@Column(nullable = true)
	private String text;

	@Enumerated(EnumType.STRING)
	@Column(nullable = false)
	private LyricsStatus lyricsStatus;

	@Column(nullable = false)
	private String syncStatus;

	@Column(nullable = false)
	private Instant createdAt;

	@Column(nullable = false)
	private Instant updatedAt;
	
	protected LyricsEntity() {}

	public LyricsEntity(
		SongEntity songEntity
	) {
		this.song = songEntity;
		this.lyricsStatus = LyricsStatus.PENDING;
		this.syncStatus = "PENDING";
	}

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
}