import { playOrPauseAudio } from "./playback.js";
import { nextSong, previousSong } from "./playlist.js";
import { getImageUrl } from "./api.js"

export function initializeMediaSession() {
  // Set metadata for the currently playing track
  if ("mediaSession" in navigator) {
    navigator.mediaSession.metadata = new MediaMetadata({
      title: "My Audio Title",
      artist: "Artist Name",
      album: "Album Name",
      artwork: [{ src: "cover.jpg", sizes: "512x512", type: "image/jpeg" }],
    });

    //   // Handle play/pause/next/previous actions
    //   navigator.mediaSession.setActionHandler('play', () => {
    //     audio.play();
    //   });

    //   navigator.mediaSession.setActionHandler('pause', () => {
    //     audio.pause();
    //   });

    navigator.mediaSession.setActionHandler("previoustrack", () => {
      playOrPauseAudio(previousSong());
    });

    navigator.mediaSession.setActionHandler("nexttrack", () => {
      playOrPauseAudio(nextSong());
    });
  }
}

export function updateMediaSessionMetadata(song, playlist) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.metadata = new MediaMetadata({
      title: playlist.name,
      artist: song.name,
      //album: playlist.name,
      artwork: [{ src: `${getImageUrl(song.thumbnail_path)}`, sizes: "512x512", type: "image/jpeg" }],
    });
  }
}
