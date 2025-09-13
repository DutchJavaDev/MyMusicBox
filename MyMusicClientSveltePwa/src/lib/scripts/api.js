const baseApiUrl = import.meta.env.VITE_BASE_API_URL;
const staticImageUrl = import.meta.env.VITE_STATIC_IMAGE_URL;
const staticAudioUrl = import.meta.env.VITE_STATIC_AUDIO_URL;

export async function fetchPlaylists() {
    try {
        const response = await fetch(`${baseApiUrl}/playlist`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const playlists = await response.json();
        return playlists.Data;
    } catch (error) {
        console.error("Error fetching playlists:", error);
        return [];
    }
}

export async function fetchNewPlaylist(lastKnowPlaylistId){
        try {
        const response = await fetch(`${baseApiUrl}/playlist?lastKnowPlaylistId=${lastKnowPlaylistId}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const playlists = await response.json();
        return playlists.Data;
    } catch (error) {
        console.error("Error fetching playlists:", error);
        return [];
    }
}

export async function fetchNewPlaylistSongs(playlistId, lastKnowSongPosition) {
        try {
        const response = await fetch(`${baseApiUrl}/playlist/${playlistId}?lastKnowSongPosition=${lastKnowSongPosition}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const songs = await response.json();
        return songs.Data;
    } catch (error) {
        console.error("Error fetching playlist songs:", error);
        return [];
    }
}

export async function fetchPlaylistSongs(playlistId) {
    try {
        const response = await fetch(`${baseApiUrl}/playlist/${playlistId}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const songs = await response.json();
        return songs.Data;
    } catch (error) {
        console.error("Error fetching playlist songs:", error);
        return [];
    }
}

export function getImageUrl(path) {
    return `${staticImageUrl}/${path}`;
}

export function getPlaybackUrl(source_id) {
    return `${staticAudioUrl}/${source_id}.opus`; // Assuming all audio files are in .opus format
}