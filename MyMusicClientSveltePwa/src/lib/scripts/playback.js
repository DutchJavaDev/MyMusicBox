let audioElement = null;

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

export let isPlaying = false;

export function initPlaybackAudio() {
    audioElement = document.getElementById("audio-player");
}

export function playAudio(playId, id) {

    isPlaying = false;

    if (!audioElement) {
        console.error("Audio element not initialized. Call initPlaybackAudio first.");
        return;
    }

    audioElement.src = `${baseApiUrl}/play/${playId}`;
    audioElement.id = id;
    audioElement.play()
    isPlaying = true;
}

