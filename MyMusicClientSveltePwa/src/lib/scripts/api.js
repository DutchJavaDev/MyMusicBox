// @ts-nocheck
import { writable } from "svelte/store";

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

let playlistsArray = new Array();
let playlistsSongsMap = new Map();

export const playlistsSongsStore = writable([]);
export const playlistsStore = writable([]);

export async function updateStores(){
    let playlists = await fetch(`${baseApiUrl}/playlist`)
        .then(response => response.json())
        .then(data => data.Data)
        .catch(error => console.error('Error fetching playlists:', error));

        playlistsArray = [];

    for (const playlist of playlists) {
        playlistsArray.push(playlist);
        let songs = await fetchPlaylistSongs(playlist.id);
        playlistsSongsMap.set(playlist.id, songs);
    }

    playlistsStore.set(playlistsArray);
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

async function fetchPlaylistSongs(playlistId) {

    let songs = await fetch(`${baseApiUrl}/playlist/${playlistId}`)
        .then(response => response.json())
        .then(data => data.Data)
        .catch(error => console.error('Error fetching playlist songs:', error));
    return songs;
}
