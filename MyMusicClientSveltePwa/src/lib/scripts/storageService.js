// @ts-ignore
const storageType = "localStorage";
const cachedPlaylistsKey = "cachedPlaylists";
const cachedPlaylistSongsKey = "cachedPlaylistSongs_";

export function getCachedPlaylists() {
  return getItem(cachedPlaylistsKey) || [];
}

export function setPlaylists(playlists) {
  setItem(cachedPlaylistsKey, playlists);
}

// Create update function to update the songs in it
export function setPlaylistSongs(playlistId, songs) {
  const key = `${cachedPlaylistSongsKey}${playlistId}`;
  setItem(key, songs);
}

export function getCachedPlaylistSongs(playlistId) {
  const key = `${cachedPlaylistSongsKey}${playlistId}`;
  return getItem(key) || [];
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
