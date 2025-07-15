import { getImageUrl } from "./api.js"
import { nextSong, previousSong } from "./playbackService.js";

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
