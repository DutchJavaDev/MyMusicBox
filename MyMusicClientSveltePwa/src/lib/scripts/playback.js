import { get, writable } from "svelte/store";
import { nextSong, getCurrentPlaylistId } from "./playlist.js";
import { updateMediaSessionMetadata } from "./mediasession.js";
import { getPlaylistById } from "./api.js";

let audioElement = null;

const baseApiUrl = import.meta.env.VITE_BASE_API_URL;

export let isPlaying = writable(false);
export let currentSong = writable(null);
export let playPercentage = writable(0);
export let isRepeatEnabled = writable(false);

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
        
        if (get(isRepeatEnabled)) {
            audioElement.currentTime = 0; // Reset to start
            audioElement.play(); // Replay the current song
        } else{
            playOrPauseAudio(nextSong());
        }
    });

    audioElement.addEventListener("timeupdate", () => {
        console.log(`Current time: ${formatTime(audioElement.currentTime)}, Duration: ${formatTime(audioElement.duration)}`);
        const percentage = (audioElement.currentTime / audioElement.duration) * 100;

        // TODO when skipping to a different song, the percentage should be reset to 0
        //Weird UI bug where the percentage is not reset to 0 when skipping to a different song
        playPercentage.set(percentage);
    })
}

function formatTime(sec) {
  const minutes = Math.floor(sec / 60);
  const seconds = Math.floor(sec % 60).toString().padStart(2, '0');
  return `${minutes}:${seconds}`;
}

export function toggleRepeat() {
    if (get(isRepeatEnabled)) {
        isRepeatEnabled.set(false);
        console.log("Repeat mode disabled");
    }
    else {
        isRepeatEnabled.set(true);
        console.log("Repeat mode enabled");
    }
}

export function setCurrentTime(seconds) {
    if (audioElement) {
        audioElement.pause();
        audioElement.currentTime = seconds;
        audioElement.play();
        console.log(`Current time set to ${seconds} seconds`);
    } else {
        console.error("Audio element is not initialized.");
    }
}

export function playOrPauseAudio(song = null) {

    if(song != null && song.source_id !== source_id) {
        playPercentage.set(0);
        audioElement.pause();
        currentSong.set(song);
        audioElement.src = `${baseApiUrl}/play/${song.source_id}`;
        source_id = song.source_id;
        audioElement.load();
        audioElement.play();
        updateMediaSessionMetadata(song, getPlaylistById(getCurrentPlaylistId()));
        return;
    }

    if(get(isPlaying)) {
        audioElement.pause();
    } else {
        audioElement.play();
    }
}
