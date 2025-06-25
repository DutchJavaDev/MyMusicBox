// @ts-ignore
const storageType = "localStorage";
const storagePrefix = "mmb_";
const playlistsKey = `${storagePrefix}playlists`;
const songsKey = `${storagePrefix}songs`;
const currentSongKey = `${storagePrefix}currentSong`;
const currentPlaylistKey = `${storagePrefix}currentPlaylist`;

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

export function storePlaylists(playlists) {
    setItem(playlistsKey, playlists);
}

export function storePlaylistSongs(playlistId, songs) {
    setItem(`${songsKey}_${playlistId}`, songs);
}

export function storeCurrentSong(index, time = 0) {
    setItem(currentSongKey, { index, time });
}

export function storeCurrentPlaylist(playlist, id, shuffle = false, repeat = false) {
    let data = { playlist, id, shuffle, repeat };
    console.log("Storing current playlist:", data);
    setItem(currentPlaylistKey, data);
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

export function getPlaylistsStore() {
  return getItem(playlistsKey) || [];
}

export function getPlaylistSongsStore(playlistId) {
  return getItem(`${songsKey}_${playlistId}`) || [];    
}

export function getCurrentSongStore() {
  return getItem(currentSongKey) || { index: 0, time: 0 };   
}

export function getCurrentPlaylist() {
  return getItem(currentPlaylistKey) || { playlist: null, shuffle: false, repeat: false };
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
