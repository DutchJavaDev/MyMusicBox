// @ts-nocheck
import { get, writable } from "svelte/store";
import { getPlaybackUrl } from "./api";
import { getCachedPlaylistSongs, getCachedPlaylists } from "./storageService";
import { shuffleArray } from "./util";
import { updateMediaSessionMetadata, updateMediaSessionPlaybackState } from "./mediasessionService";

export let currentSong = writable(null);
export let isPlaying = writable(false);
export let playPercentage = writable(0);
export let isShuffledEnabled = writable(false);
export let isLoopingEnabled = writable(false);

let audioElement = null;
let currentPlaylistId = 0;
let originalPlaylistSongs = [];
let playlistSongs = [];
let songIndex = 0;

export function initializePlaybackService() {
  audioElement = document.getElementById("audio-player");
  if (!audioElement) {
    console.error("Audio element with id 'audio-player' not found in the document.");
    return;
  }

  audioElement.addEventListener("play", () => {
    isPlaying.set(true);
    updateMediaSessionPlaybackState(true);
    const currentPlaylist = getCachedPlaylists().find(pl => pl.id === currentPlaylistId);
    updateMediaSessionMetadata(get(currentSong), currentPlaylist);
  });

  audioElement.addEventListener("ended", () => {
    isPlaying.set(false);
    updateMediaSessionPlaybackState(false);
    if (get(isLoopingEnabled)) {
      audioElement.currentTime = 0;
      audioElement.play();
    } else {
      nextSong();
    }
  });

  audioElement.addEventListener("pause", () => {
    isPlaying.set(false);
    updateMediaSessionPlaybackState(false);
  });

  audioElement.addEventListener("playing", () => {
    isPlaying.set(true);
    updateMediaSessionPlaybackState(true);
  });

  audioElement.addEventListener("timeupdate", () => {
    const percentage = (audioElement.currentTime / audioElement.duration) * 100;

    // Fixes ui bug
    if (isNaN(percentage)) {
      return;
    }

    playPercentage.set(percentage);
  });
}

export function nextSong() {
  songIndex = (songIndex + 1) % playlistSongs.length; // Loop back to start if at end of playlist
  const nextSong = playlistSongs[songIndex];
  playOrPauseSong(nextSong.id);
}

export function previousSong() {
  songIndex = (songIndex - 1 + playlistSongs.length) % playlistSongs.length; // Loop to end if at start of playlist
  const previousSong = playlistSongs[songIndex];
  playOrPauseSong(previousSong.id);
}

export function playOrPauseSong(songId) {
  const _currentSong = get(currentSong);

  if (!_currentSong || _currentSong.id != songId) {
    // new song
    let song = playlistSongs.find((song) => song.id === songId);
    songIndex = playlistSongs.findIndex((song) => song.id === songId);
    audioElement.pause();
    audioElement.src = getPlaybackUrl(song.source_id);
    audioElement.load();
    currentSong.set(playlistSongs.find((song) => song.id === songId));
    isPlaying.set(false); // set to false since this is a new song
  }

  if (get(isPlaying)) {
    audioElement.pause();
  } else {
    audioElement.play(); // https://developer.chrome.com/blog/play-request-was-interrupted
  }
}

export function toggleShuffle() {
  if (get(isShuffledEnabled)) {
    playlistSongs = originalPlaylistSongs;
    songIndex = playlistSongs.findIndex((song) => song.id === get(currentSong).id);
    isShuffledEnabled.set(false);
  } else {
    playlistSongs = shuffleArray([...originalPlaylistSongs]);
    songIndex = playlistSongs.findIndex((song) => song.id === get(currentSong).id);
    isShuffledEnabled.set(true);
  }
}

export function toggleLoop() {
    isLoopingEnabled.set(!get(isLoopingEnabled));
}

export function setCurrentTime(seconds) {
  audioElement.pause();
  audioElement.currentTime = seconds;
  audioElement.play();
}

export function setPlaylists(playlistId) {
  if (currentPlaylistId === playlistId) {
    // update current playlist
    originalPlaylistSongs = getCachedPlaylistSongs(playlistId);
    return; // Already set to this playlist
  }
  currentPlaylistId = playlistId;
  originalPlaylistSongs = getCachedPlaylistSongs(playlistId);
  playlistSongs = originalPlaylistSongs;
}
