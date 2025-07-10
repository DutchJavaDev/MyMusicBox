// @ts-nocheck
import { writable, get } from "svelte/store";
import { storePlaylists, storePlaylistSongs, getPlaylistsStore, getPlaylistSongsStore } from "./storage.js";

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

const playlistsArray = new Array();
const playlistsSongsMap = new Map();

export const playlistsStore = writable([]);

export const writablePlaylistsStore = [];

export async function initStores() {
  // check localStorage for playlists and songs
  // if not found, fetch from API
  const cachedPlaylists = getPlaylistsStore();

  if (cachedPlaylists.length > 0) {
    for (const playlist of cachedPlaylists) {
      playlistsArray.push(playlist);
      let songs = getPlaylistSongsStore(playlist.id);
      playlistsSongsMap.set(playlist.id, songs);

      writablePlaylistsStore[playlist.id] = writable(songs);
    }

    playlistsStore.set(playlistsArray);
    console.log("Loaded playlists from localStorage");
  } else {
    console.log("Fetching playlists from API");
    const playlists = await fetch(`${baseApiUrl}/playlist`)
      .then((response) => response.json())
      .then((data) => data.Data)
      .catch((error) => console.error("Error fetching playlists:", error));

    playlistsArray.length = 0;;
    playlistsSongsMap.clear();

    for (const playlist of playlists) {
      playlistsArray.push(playlist);
      const songs = await fetchPlaylistSongs(playlist.id);
      playlistsSongsMap.set(playlist.id, songs);

      storePlaylistSongs(playlist.id, songs);

      writablePlaylistsStore[playlist.id] = writable(songs);
    }

    playlistsStore.set(playlistsArray);
    storePlaylists(playlistsArray);
  }
}

export async function updateStores() {
  // Update plaists
  const cachedPlaylists = getPlaylistsStore();
  const lastKnowPlaylistId = cachedPlaylists.at(-1).id;
  const playlists = await fetchPlaylists(lastKnowPlaylistId);

  for (const playlist of playlists) {
    playlistsArray.push(playlist);
  }

  playlistsStore.set(playlistsArray);
  storePlaylists(playlistsArray);

  // Update songs for each playlist
  console.log("starting to update songs for each playlist");
  for (const playlist of playlistsArray) {
    
    let lastKnowSongPosition = playlistsSongsMap.get(playlist.id).length;

    if (!lastKnowSongPosition) {
      lastKnowSongPosition = 0; // Default to 0 if no songs are found
    }

    const songs = await fetchPlaylistSongs(playlist.id, lastKnowSongPosition);

    if (songs.length > 0) {
      // Add local notification for new songs?
      console.log(`Found (${songs.length}) new songs for playlist ID: ${playlist.id} with last known song position: ${lastKnowSongPosition}`);
      playlistsSongsMap.set(playlist.id, [...playlistsSongsMap.get(playlist.id), ...songs]);
      storePlaylistSongs(playlist.id, playlistsSongsMap.get(playlist.id));

      const songList = get(writablePlaylistsStore[playlist.id]);

      for (const song of songs) {
        songList.push(song);
      }

      writablePlaylistsStore[playlist.id].set(songList);
    }
  }
  console.log("Finished updating songs for each playlist");
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

async function fetchPlaylists(lastKnowPlaylistId) {
  const playlists = await fetch(`${baseApiUrl}/playlist?lastKnowPlaylistId=${lastKnowPlaylistId}`)
    .then((response) => response.json())
    .then((data) => data.Data)
    .catch((error) => console.error("Error fetching playlists:", error));
  return playlists;
}

async function fetchPlaylistSongs(playlistId, lastKnowSongId = 0) {
  const songs = await fetch(`${baseApiUrl}/playlist/${playlistId}?lastKnowSongPosition=${lastKnowSongId}`)
    .then((response) => response.json())
    .then((data) => data.Data)
    .catch((error) => console.error("Error fetching playlist songs:", error));
  return songs;
}
