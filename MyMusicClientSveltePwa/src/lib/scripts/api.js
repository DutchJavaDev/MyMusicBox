// @ts-nocheck
import { writable } from "svelte/store";
import { storePlaylists, storePlaylistSongs, getPlaylistsStore, getPlaylistSongsStore } from "./storage.js";

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

let playlistsArray = new Array();
let playlistsSongsMap = new Map();

export const playlistsStore = writable([]);

export async function updateStores() {
  // check localStorage for playlists and songs
  // if not found, fetch from API
  let cachedPlaylists = getPlaylistsStore();

  if (cachedPlaylists.length > 0) {
    for (const playlist of cachedPlaylists) {
      playlistsArray.push(playlist);  
      let songs = getPlaylistSongsStore(playlist.id);
      playlistsSongsMap.set(playlist.id, songs);
    }

    playlistsStore.set(playlistsArray);
    console.log("Loaded playlists from localStorage");
  } else {
    let playlists = await fetch(`${baseApiUrl}/playlist`)
      .then((response) => response.json())
      .then((data) => data.Data)
      .catch((error) => console.error("Error fetching playlists:", error));

    playlistsArray = [];
    playlistsSongsMap.clear();

    for (const playlist of playlists) {
      playlistsArray.push(playlist);
      let songs = await fetchPlaylistSongs(playlist.id);
      playlistsSongsMap.set(playlist.id, songs);

      storePlaylistSongs(playlist.id, songs);
    }

    playlistsStore.set(playlistsArray);
    storePlaylists(playlistsArray);
  }
}

export function getImageUrl(imagePath) {
  if (!imagePath) return null;
  return `${baseApiUrl}/images/${imagePath}`;
}

export function getPlaylistSongs(playlistId) {
  if (playlistsSongsMap.has(playlistId)) {
    return playlistsSongsMap.get(playlistId);
  } else {
    console.warn(`No songs found for playlist ID: ${playlistId}`);
    return [];
  }
}

export function getPlaylistById(playlistId) {
  const playlist = playlistsArray.find((p) => p.id === playlistId);
  if (playlist) {
    return playlist;
  } else {
    console.warn(`No playlist found with ID: ${playlistId}`);
    return null;
  }
}

async function fetchPlaylistSongs(playlistId) {
  let songs = await fetch(`${baseApiUrl}/playlist/${playlistId}`)
    .then((response) => response.json())
    .then((data) => data.Data)
    .catch((error) => console.error("Error fetching playlist songs:", error));
  return songs;
}
