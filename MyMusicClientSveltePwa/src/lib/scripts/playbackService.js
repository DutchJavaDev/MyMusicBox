// @ts-nocheck
import { get, writable } from "svelte/store";
import { getPlaybackUrl } from "./api";
import { getCachedPlaylistSongs, getCachedPlaylists, 
         setPlaybackState, getPlaybackState,
         getCurrentPlaylistId, setCurrentPlaylistId,
         getCurrentShuffledPlaylist, setCurrentShuffledPlaylist,
         getCurrentSongIndex, setCurrentSongIndex, setCurrentSongTime,
         getCurrentSongTime } from "./storageService";
import { shuffleArray } from "./util";
import { updateMediaSessionMetadata, updateMediaSessionPlaybackState, updatePositionState } from "./mediasessionService";

export let currentSong = writable({id: -999, title: "", artist: "", album: "", source_id: ""});
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

  const playbackState = getPlaybackState();
  isLoopingEnabled.set(playbackState.isLoopingEnabled);
  isShuffledEnabled.set(playbackState.isShuffledEnabled);

  // if playbackState.isShuffledEnabled is true, we need to get the shuffled playlist
  if(get(isShuffledEnabled)) {
    currentPlaylistId = getCurrentPlaylistId();
    playlistSongs = getCurrentShuffledPlaylist();
    originalPlaylistSongs = getCachedPlaylistSongs(currentPlaylistId);
    songIndex = getCurrentSongIndex();
    isPlaying.set(false);
    setCurrentTime(getCurrentSongTime());
    playOrPauseSong(playlistSongs[songIndex].id);
  } else{
    currentPlaylistId = getCurrentPlaylistId();
    originalPlaylistSongs = getCachedPlaylistSongs(currentPlaylistId);
    playlistSongs = originalPlaylistSongs;
    songIndex = getCurrentSongIndex();
    isPlaying.set(false);
    setCurrentTime(getCurrentSongTime());
    playOrPauseSong(playlistSongs[songIndex].id);
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
      audioElement.load();
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

    setCurrentSongTime(audioElement.currentTime);
    updatePositionState(audioElement.currentTime, audioElement.duration);
    playPercentage.set(percentage);
  });

  audioElement.addEventListener("loadeddata", async (e) => {
    // Safe to play audio now
    await audioElement.play();
  });
  audioElement.addEventListener("error", (e) => {
    console.error("Error loading audio:", e);
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
    // new song selected
    playPercentage.set(0);
    let song = playlistSongs.find((song) => song.id === songId);
    songIndex = playlistSongs.findIndex((song) => song.id === songId);
    audioElement.src = getPlaybackUrl(song.source_id);
    audioElement.load();
    currentSong.set(playlistSongs.find((song) => song.id === songId));
    isPlaying.set(false); // set to false since this is a new song
    setCurrentSongIndex(songIndex);
  }
  else if (get(isPlaying)) {
    audioElement.pause();
  }else {
    // data is already loaded, just play
    audioElement.play();
  }
}

export function toggleShuffle() {
  if (get(isShuffledEnabled)) {
    playlistSongs = originalPlaylistSongs;
    songIndex = playlistSongs.findIndex((song) => song.id === get(currentSong).id);
    isShuffledEnabled.set(false);

    setCurrentShuffledPlaylist([]);
    setCurrentSongIndex(songIndex);
  } else {
    playlistSongs = shuffleArray([...originalPlaylistSongs]);
    songIndex = playlistSongs.findIndex((song) => song.id === get(currentSong).id);
    isShuffledEnabled.set(true);

    setCurrentShuffledPlaylist(playlistSongs);
    setCurrentSongIndex(songIndex);
  }

  setPlaybackState(get(isLoopingEnabled), get(isShuffledEnabled));
}

export function toggleLoop() {
    isLoopingEnabled.set(!get(isLoopingEnabled));
    setPlaybackState(get(isLoopingEnabled), get(isShuffledEnabled));
}

export function setCurrentTime(seconds) {
  audioElement.currentTime = seconds;
}

export function updateCurrentPlaylist(playlistId){
if (currentPlaylistId === playlistId) {
    // update current playlist
    originalPlaylistSongs = getCachedPlaylistSongs(playlistId);

    // Todo get the difference between originalPlaylistSongs and playlistSongs
    // and update playlistSongs accordingly if shuffle is enabled
    return; // Already set to this playlist
  }
}

export function setPlaylists(playlistId) {
  if (currentPlaylistId === playlistId) {
    // update current playlist
    originalPlaylistSongs = getCachedPlaylistSongs(playlistId);

    // Todo get the difference between originalPlaylistSongs and playlistSongs
    // and update playlistSongs accordingly if shuffle is enabled
    return; // Already set to this playlist
  }
  isLoopingEnabled.set(false);
  isShuffledEnabled.set(false);
  setPlaybackState(false, false);
  currentPlaylistId = playlistId;
  originalPlaylistSongs = getCachedPlaylistSongs(playlistId);
  playlistSongs = originalPlaylistSongs;
  setCurrentPlaylistId(currentPlaylistId);
}
