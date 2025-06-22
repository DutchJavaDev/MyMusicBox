import { playOrPauseAudio } from "./playback.js";
import { nextSong, previousSong } from "./playlist.js";
import { getImageUrl } from "./api.js"

export function initializeMediaSession() {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.setActionHandler("previoustrack", () => {
      playOrPauseAudio(previousSong());
    });

    navigator.mediaSession.setActionHandler("nexttrack", () => {
      playOrPauseAudio(nextSong());
    });
  }
}

export function updateMediaSessionPlaybackState(isPlaying) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.playbackState = isPlaying ? "playing" : "paused";
  }
}

export function updateMediaSessionMetadata(song, playlist) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.playbackState = "playing";
    navigator.mediaSession.metadata = new MediaMetadata({
      title: playlist.name,
      artist: song.name,
      artwork: [{ src: `${getImageUrl(song.thumbnail_path)}`, sizes: "512x512", type: "image/jpeg" }],
    });
  }
}
