import { writable, get } from "svelte/store";
import { getCachedPlaylists, setCachedPlaylists, setCachedPlaylistSongs, getCachedPlaylistSongs, appConfiguration, getConfiguration, getCurrentPlaylistId } from "./storageService";
import { fetchPlaylists, fetchPlaylistSongs, fetchNewPlaylist, fetchNewPlaylistSongs, deletePlaylist } from "./api";
import { componentParams, navigateTo } from "./routeService";
import { playOrPauseSong, setPlaylists, currentSong, playPercentage, updateCurrentPlaylist } from "./playbackService";

export const playlistsStore = writable([]);

let _componentParams;

let updateInterval;
let isUpdating = false;
let intervalId;

// Check storage for stored playlists, if empty fetch from API
export async function initializePlaylistService() {
  componentParams.subscribe((value) => {
    _componentParams = value;
  });

  const cachedPlaylists = getCachedPlaylists();
  if (cachedPlaylists.length > 0) {
    playlistsStore.set(cachedPlaylists);
  } else {
    const fetchedPlaylists = await fetchPlaylists();
    playlistsStore.set(fetchedPlaylists);
    setCachedPlaylists(fetchedPlaylists);
    for (const playlist of fetchedPlaylists) {
      const songs = await fetchPlaylistSongs(playlist.id);
      setCachedPlaylistSongs(playlist.id, songs);
    }
  }

  updateInterval = getConfiguration().fetchTimer * 1000; // Need to multiply by 1000 to get milliseconds
  // Subscribe to configuration changes
  // If fetchTimer is updated, clear the old interval and set a new one

  appConfiguration.subscribe((config) => {
    if (intervalId) {
      clearInterval(intervalId);
    }
    updateInterval = config.fetchTimer * 1000; // Need to multiply by 1000 to get milliseconds

    intervalId = setInterval(() => {
      if (isUpdating) return; // Prevent multiple updates at the same time

      isUpdating = true;

      backgroundUpdate();

      isUpdating = false;
    }, updateInterval);
  });
}

export async function deleteCurrentPlaylist() {
  const playlistId = _componentParams.playlistId;

  // check if playlistId actually exists in cached playlists
  const cachedPlaylists = getCachedPlaylists();
  const playlistIndex = cachedPlaylists.findIndex((/** @type {{ id: any; }} */ p) => p.id === playlistId);

  if (playlistIndex === -1) {
    alert("Playlist not found.");
    return;
  }

  if (confirm("Are you sure you want to delete the current playlist? This action cannot be undone.")) {
    // TODO delete resource from API, images etc
    const result = await deletePlaylist(playlistId);
    if (result.success) {
      const currentPlaylistId = getCurrentPlaylistId();

      // If the deleted playlist is the current playing playlist, stop playback
      if (currentPlaylistId === playlistId) {
        // stop playback
        playOrPauseSong(null);
        setPlaylists(0);
        currentSong.set({ id: -999, title: "", artist: "", album: "", source_id: "" });
      }

      playPercentage.set(0);

      // Remove playlist from cached playlists
      cachedPlaylists.splice(playlistIndex, 1);
      setCachedPlaylists(cachedPlaylists);
      playlistsStore.set(cachedPlaylists);
      alert("Playlist deleted successfully.");

      navigateTo("/");
    }
  }
}

export function removeSongFromPlaylist(songId) {
  const cachedPlaylists = getCachedPlaylists();

  // remove song from playlist
  for (const playlist of cachedPlaylists) {
    const playlistId = playlist.id;
    const cachedSongs = getCachedPlaylistSongs(playlistId);
    const songIndex = cachedSongs.findIndex((s) => s.id === songId);
    const removed = cachedSongs.splice(songIndex, 1);

    if (removed.length > 0) {
      setCachedPlaylistSongs(playlistId, cachedSongs);
    }
  }
}

async function backgroundUpdate() {
  // update playlists in the background
  const lastItemInex = -1;
  const cachedPlaylists = getCachedPlaylists();
  const lastKnowPlaylistId = cachedPlaylists.at(lastItemInex).id;

  const newPlaylists = await fetchNewPlaylist(lastKnowPlaylistId);

  if (newPlaylists.length > 0) {
    cachedPlaylists.push(...newPlaylists);
    setCachedPlaylists(cachedPlaylists);
    playlistsStore.set(cachedPlaylists);
  }

  //update songs in the background
  for (const playlist of cachedPlaylists) {
    const playlistId = playlist.id;
    const cachedSongs = getCachedPlaylistSongs(playlistId);
    const lastKnowSongPosition = cachedSongs.length;

    const newPlaylistSongs = await fetchNewPlaylistSongs(playlistId, lastKnowSongPosition);
    
    // get the difference between cachedSongs and newPlaylistSongs
    const songsToAdd = newPlaylistSongs.filter(song => !cachedSongs.some(cs => cs.id === song.id));

    if (songsToAdd.length > 0) {
      const updatedSongs = [...cachedSongs, ...songsToAdd];
      setCachedPlaylistSongs(playlistId, updatedSongs);
    } 

  }
}
