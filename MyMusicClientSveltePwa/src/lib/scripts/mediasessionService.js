import { getImageUrl } from "./api.js"
import { nextSong, previousSong } from "./playbackService.js";

const DEFAULT_PLAYBACK_RATE = 1.0;

export function initializeMediaSessionService() {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.setActionHandler("previoustrack", () => {
      previousSong();
    });

    navigator.mediaSession.setActionHandler("nexttrack", () => {
      nextSong();
    });
  }
}

export function updateMediaSessionPlaybackState(isPlaying) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.playbackState = "none";
    navigator.mediaSession.playbackState = isPlaying ? "playing" : "paused";
  }
}

export function updatePositionState(currentTime, duration) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.setPositionState({
      duration: duration,
      playbackRate: DEFAULT_PLAYBACK_RATE,
      position: currentTime,
    });
  }
}

export function updateMediaSessionMetadata(song, playlist) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.playbackState = "paused";
    navigator.mediaSession.playbackState = "playing";
    navigator.mediaSession.metadata = new MediaMetadata({
      title: playlist.name,
      artist: song.name,
      album: playlist.name,
      artwork: [{ src: `${getImageUrl(song.thumbnail_path)}`, sizes: "512x512", type: "image/jpeg" }],
    });
  }
}
