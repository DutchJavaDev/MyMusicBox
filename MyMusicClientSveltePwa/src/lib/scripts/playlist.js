import { get, writable } from "svelte/store";
import { getPlaylistSongs } from "./api.js";

export let currentPlaylistId = writable(-1);
export let isShuffleEnabled = writable(false);
let currentIndex = 0;
let originalPlaylist = []; // This can be used to store the original playlist if needed
let currentPlaylist = [];

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
}

export function getCurrentSong() {
  return currentPlaylist[currentIndex];
}

export function setCurrentSong(song) {
  const index = currentPlaylist.findIndex((s) => s.id === song.id);
  if (index !== -1) {
    currentIndex = index; // Update the current index to the song's index
  } else {
    console.warn("Song not found in the current playlist.");
  }
}

export function nextSong() {
  if (currentPlaylist.length === 0) return null; // No songs to play
  currentIndex = (currentIndex + 1) % currentPlaylist.length; // Loop back to start
  return getCurrentSong();
}

export function previousSong() {
  if (currentPlaylist.length === 0) return null; // No songs to play
  currentIndex = (currentIndex - 1 + currentPlaylist.length) % currentPlaylist.length; // Loop back to end
  return getCurrentSong();
}

export function shufflePlaylist() {
  if (get(isShuffleEnabled)) {
    currentPlaylist = originalPlaylist.slice();
    currentIndex = 0;
    isShuffleEnabled.set(false);
  } else {
    currentPlaylist = currentPlaylist.sort(() => Math.random() - 0.5);
    currentIndex = 0;
    isShuffleEnabled.set(true);
  }
}
