import { get, writable } from "svelte/store";
import { nextSong } from "./playlist.js";

let audioElement = null;

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

export let isPlaying = writable(false);
export let currentSong = writable(null);
export let playPercentage = writable(0);

let source_id;

export function initPlaybackAudio() {
    audioElement = document.getElementById("audio-player");
    if (!audioElement) {
        console.error("Audio element with id 'audio-player' not found in the document.");
        return;
    }

    audioElement.addEventListener("playing", () => {
        console.log("Audio is playing");
        isPlaying.set(true);    
    });

    audioElement.addEventListener("pause", () => {
        console.log("Audio is paused");
        isPlaying.set(false);
    });

    audioElement.addEventListener("ended", () => {
        console.log("Audio has ended");
        isPlaying.set(false);
        playPercentage.set(0);
        playOrPauseAudio(nextSong());
    });

    audioElement.addEventListener("timeupdate", () => {
        console.log(`Current time: ${formatTime(audioElement.currentTime)}, Duration: ${formatTime(audioElement.duration)}`);
        const percentage = (audioElement.currentTime / audioElement.duration) * 100;
        playPercentage.set(percentage);
    })
}

function formatTime(sec) {
  const minutes = Math.floor(sec / 60);
  const seconds = Math.floor(sec % 60).toString().padStart(2, '0');
  return `${minutes}:${seconds}`;
}

export function playOrPauseAudio(song = null) {

    if(song != null && song.source_id !== source_id) {
        audioElement.pause();
        currentSong.set(song);
        audioElement.src = `${baseApiUrl}/play/${song.source_id}`;
        source_id = song.source_id;
        audioElement.load();
        audioElement.play();
        return;
    }

    if(get(isPlaying)) {
        audioElement.pause();
    } else {
        audioElement.play();
    }
}
