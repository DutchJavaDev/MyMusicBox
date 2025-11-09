import { get, writable } from "svelte/store";

// @ts-ignore
const storageType = "localStorage";
const PlaylistsKey = "cachedPlaylists";
const PlaylistSongsKey = "cachedPlaylistSongs_";
const PlaybackStateKey = "playbackState";
const CurrentPlaylistIdKey = "currentPlaylistId";
const CurrentSongIndexKey = "currentSongIndex";
const CurrentShuffeldPlaylistKey = "currentShuffledPlaylist";
const CurrentSongTimeKey = "currentSongTime";
const ConfigKey = "appConfig";
const ConfigKeySearchQyery = "searchQuery";

export let appConfiguration = writable(getConfiguration());

export function setConfiguration(config) {
  appConfiguration.set(config);
  setItem(ConfigKey, config);
}

export function setSearchQueryInStorage(query) {
  setItem(ConfigKeySearchQyery, query);
} 

export function getSearchQueryFromStorage() {
  return getItem(ConfigKeySearchQyery) || '';
} 

export function setCachedPlaylists(playlists) {
  setItem(PlaylistsKey, playlists);
}

export function setCachedPlaylistSongs(playlistId, songs) {
  const key = `${PlaylistSongsKey}${playlistId}`;
  setItem(key, songs);
}

export function setPlaybackState(isLoopingEnabled, isShuffledEnabled) {
  const playbackState = {
    isLoopingEnabled,
    isShuffledEnabled,
  };
  setItem(PlaybackStateKey, playbackState);
}

export function setCurrentPlaylistId(playlistId) {
  setItem(CurrentPlaylistIdKey, playlistId);
}

export function setCurrentSongIndex(index) {
  setItem(CurrentSongIndexKey, index);
}

export function setCurrentShuffledPlaylist(shuffledPlaylist) {
  setItem(CurrentShuffeldPlaylistKey, shuffledPlaylist);
}

export function setCurrentSongTime(seconds) {
  setItem(CurrentSongTimeKey, seconds);
}

export function getCachedPlaylists() {
  return getItem(PlaylistsKey) || [];
}

export function getCachedPlaylistSongs(playlistId) {
  const key = `${PlaylistSongsKey}${playlistId}`;
  return getItem(key) || [];
}

export function getPlaybackState() {
  return getItem(PlaybackStateKey) || {
    isLoopingEnabled: false,
    isShuffledEnabled: false,
  };
}

export function getCurrentPlaylistId() {
  return getItem(CurrentPlaylistIdKey) || 0;
}

export function getCurrentSongIndex() {
  return getItem(CurrentSongIndexKey) || 0;
}

export function getCurrentShuffledPlaylist() {
  return getItem(CurrentShuffeldPlaylistKey) || [];
}

export function getCurrentSongTime() {
  return getItem(CurrentSongTimeKey) || 0;
}

export function getConfiguration() {
  return getItem(ConfigKey) || { sleepTimer: 15, // minutes 
                                 fetchTimer: 3,  // seconds
                                 byPassCache: false 
                                };
}

export function clearStorage() {
    if (storageAvailable(storageType)) {
        try {
            localStorage.clear();
            console.log("Local storage cleared.");
        } catch (e) {
            console.error("Error clearing localStorage:", e);
        }
    }
}

function setItem(key, value) {
  if (storageAvailable(storageType)) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch (e) {
      console.error("Error setting item in localStorage:", e);
    }
  }
}

function getItem(key) {
  if (storageAvailable(storageType)) {
    try {
      const value = localStorage.getItem(key);
      return value ? JSON.parse(value) : null;
    } catch (e) {
      console.error("Error getting item from localStorage:", e);
      return null;
    }
  }
  return null;
}

function storageAvailable(type) {

  let storage;
  try {
    storage = window[type];
    const x = "__storage_test__";
    // @ts-ignore
    storage.setItem(x, x);
    // @ts-ignore
    storage.removeItem(x);
    return true;
  } catch (e) {
    alert(`Storage Error: ${e.message}`);
    return (
      e instanceof DOMException &&
      e.name === "QuotaExceededError" &&
      // acknowledge QuotaExceededError only if there's something already stored
      storage &&
      storage.length !== 0
    );
  }
}
