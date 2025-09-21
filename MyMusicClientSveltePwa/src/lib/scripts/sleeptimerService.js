// @ts-nocheck
import { writable, get } from "svelte/store";
import { getConfiguration } from "./storageService";

export let timeLeft = writable(0);
export let isTimerEnabled = writable(false);

const Seconds = 1000;
const Minutes = 60;

let timeOutId;
let intervalId;

export function toggleSleepTimer() {
    
if (get(isTimerEnabled)) {
  clearTimeout(timeOutId);
  clearInterval(intervalId);
  timeOutId = null;
  console.log("Sleep timer cleared.");
  isTimerEnabled.set(false);
  return;
}

const config = getConfiguration();

const totalMinutes = config.sleepTimer; // Default to 30 minutes if no time is provided
let remainingMinutes = totalMinutes;

timeLeft.set(remainingMinutes);
isTimerEnabled.set(true);

// Start the timeout
timeOutId = setTimeout(() => {
    let audioElement = document.getElementById("audio-player");
  if (!audioElement) {
    console.error("Audio element with id 'audio-player' not found in the document.");
    return;
  }
    audioElement.pause();
}, totalMinutes * Minutes * Seconds);

// Display countdown every minute
intervalId = setInterval(() => {
  remainingMinutes--;

  if (remainingMinutes > 0) {
    timeLeft.set(remainingMinutes);
  } else {
    clearInterval(intervalId);
    isTimerEnabled.set(false);
    clearTimeout(timeOutId);
    timeOutId = null;
  }
}, Minutes * Seconds); // Every minute

}