import { writable } from "svelte/store";
import { getCachedPlaylists, setPlaylists, setPlaylistSongs, getCachedPlaylistSongs, appConfiguration, getConfiguration } from "./storageService";
import { fetchPlaylists, fetchPlaylistSongs, fetchNewPlaylist, fetchNewPlaylistSongs } from "./api";

export const playlistsStore = writable([]);

let updateInterval; 
let isUpdating = false;
let intervalId;

// Check storage for stored playlists, if empty fetch from API
export async function initializePlaylistService() {
  const cachedPlaylists = getCachedPlaylists();
  if (cachedPlaylists.length > 0) {
    playlistsStore.set(cachedPlaylists);
  } else {
    const fetchedPlaylists = await fetchPlaylists();
    playlistsStore.set(fetchedPlaylists);
    setPlaylists(fetchedPlaylists);
    for (const playlist of fetchedPlaylists) {
      const songs = await fetchPlaylistSongs(playlist.id);
      setPlaylistSongs(playlist.id, songs);
    }
  }

  updateInterval = getConfiguration().fetchTimer * 1000; // Need to multiply by 1000 to get milliseconds
  // Subscribe to configuration changes
  // If fetchTimer is updated, clear the old interval and set a new one
  
  appConfiguration.subscribe(config => {
    if (intervalId) {
      clearInterval(intervalId);
    }
    updateInterval = config.fetchTimer * 1000; // Need to multiply by 1000 to get milliseconds

    console.log("Update interval set to:", updateInterval, "ms");

    intervalId = setInterval(() => {
    
    if (isUpdating) return; // Prevent multiple updates at the same time
    
    isUpdating = true;
    
    backgroundUpdate();
    
    isUpdating = false;

  }, updateInterval);
  });
}

async function backgroundUpdate() {
  // update playlists in the background
  const lastItemInex = -1;
  const cachedPlaylists = getCachedPlaylists();
  const lastKnowPlaylistId = cachedPlaylists.at(lastItemInex).id;

  const newPlaylists = await fetchNewPlaylist(lastKnowPlaylistId);

  if (newPlaylists.length > 0) {
    cachedPlaylists.push(...newPlaylists);
    setPlaylists(cachedPlaylists);
    playlistsStore.set(cachedPlaylists);
  }

  //update songs in the background
  for (const playlist of cachedPlaylists) {
    const playlistId = playlist.id;
    const cachedSongs = getCachedPlaylistSongs(playlistId);
    const lastKnowSongPosition = cachedSongs.length;
    const newSongs = await fetchNewPlaylistSongs(playlistId, lastKnowSongPosition);
    if (newSongs.length > 0) {
      const updatedSongs = [...cachedSongs, ...newSongs];
      setPlaylistSongs(playlistId, updatedSongs);
    }
  }
}

