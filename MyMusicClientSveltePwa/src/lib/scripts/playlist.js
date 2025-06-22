import { get, writable } from "svelte/store";
import { getPlaylistSongs } from "./api.js";
import { storeCurrentPlaylist, storeCurrentSong, getCurrentPlaylist, getCurrentSongStore, getPlaylistSongsStore } from "./storage.js";
import { isRepeatEnabled, toggleRepeat } from "./playback.js";

export let currentPlaylistId = writable(-1);
export let isShuffleEnabled = writable(false);
let currentIndex = 0;
let originalPlaylist = []; // This can be used to store the original playlist if needed
let currentPlaylist = [];

export function initPlaylist() {
  const currentPlaylistData = getCurrentPlaylist();
  if (currentPlaylistData.playlist) {
    currentPlaylistId.set(currentPlaylistData.id);
    originalPlaylist = getPlaylistSongsStore(currentPlaylistData.playlist.id);
    currentPlaylist = currentPlaylistData.playlist;
    isShuffleEnabled.set(currentPlaylistData.shuffle);
    isRepeatEnabled.set(currentPlaylistData.repeat);

    const currentSong = getCurrentSongStore();
    if (currentSong && currentSong.index !== undefined) {
      currentIndex = currentSong.index;
    }
  } else {
    console.warn("No current playlist found in storage.");
  }
}

export function getCurrentPlaylistId() {
  return get(currentPlaylistId);
}

export function setSongs(playlistId) {
  if (get(currentPlaylistId) === playlistId) {
    return; // No need to update if the same playlist is set
  }
  currentPlaylistId.set(playlistId); // Update the current playlist ID
  originalPlaylist = getPlaylistSongs(playlistId);
  currentPlaylist = originalPlaylist.slice(); // Create a copy of the original playlist
  currentIndex = 0; // Reset index when setting new songs

  storeCurrentPlaylist(currentPlaylist, playlistId, get(isShuffleEnabled), get(isRepeatEnabled));
}

export function getCurrentSong() {
  let song = currentPlaylist[currentIndex];
  storeCurrentSong(currentIndex, 0);
  return song;
}

export function getCurrentIndex() {
  return currentIndex;
}

export function setCurrentSong(song) {
  const index = currentPlaylist.findIndex((s) => s.id === song.id);
  if (index !== -1) {
    currentIndex = index; // Update the current index to the song's index
  } else {
    console.warn("Song not found in the current playlist.");
  }
  storeCurrentSong(currentIndex, 0); // Store the current song with time reset
}

export function nextSong() {
  if (currentPlaylist.length === 0) return null; // No songs to play
  currentIndex = (currentIndex + 1) % currentPlaylist.length; // Loop back to start

  if (get(isRepeatEnabled)) {
    toggleRepeat(); // Disable repeat if it was enabled
  }

  return getCurrentSong();
}

export function previousSong() {
  if (currentPlaylist.length === 0) return null; // No songs to play
  currentIndex = (currentIndex - 1 + currentPlaylist.length) % currentPlaylist.length; // Loop back to end

  if (get(isRepeatEnabled)) {
    toggleRepeat(); // Disable repeat if it was enabled
  }

  return getCurrentSong();
}

export function shufflePlaylist() {
  const currentSong = getCurrentSong();

  if (get(isShuffleEnabled)) {
    currentPlaylist = originalPlaylist.slice();
    currentIndex = currentPlaylist.findIndex((s) => s.id === currentSong.id);
    isShuffleEnabled.set(false);
  } else {
    currentPlaylist = currentPlaylist.sort(() => Math.random() - 0.5);
    currentIndex = currentPlaylist.findIndex((s) => s.id === currentSong.id);
    isShuffleEnabled.set(true);
  }
  storeCurrentPlaylist(currentPlaylist, get(currentPlaylistId), get(isShuffleEnabled), get(isRepeatEnabled));
}
