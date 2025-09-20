import { writable } from "svelte/store";

export let searchQuery = writable();

// Mulberry32 PRNG
function mulberry32(seed) {
  return function () {
    let t = (seed += 0x6d2b79f5);
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

// Generate a stronger random 32-bit seed
function generateSeed() {
  if (typeof crypto !== "undefined" && crypto.getRandomValues) {
    const array = new Uint32Array(1);
    crypto.getRandomValues(array);
    return array[0];
  }
  // Fallback for Node.js or older browsers
  return (Date.now() ^ Math.floor(Math.random() * 0xffffffff)) >>> 0;
}

// Fisher-Yates shuffle with optional seed
export function shuffleArray(array) {
  const seed = generateSeed();
  const rng = mulberry32(seed);

  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(rng() * (i + 1));
    [array[i], array[j]] = [array[j], array[i]];
  }

  return array;
}

export function getSearchParameters() {
  const searchParams = new URLSearchParams(window.location.search);
  const result = {};
  for (const [key, value] of searchParams.entries()) {
    if (result[key]) {
      // If key already exists, convert to array or push into existing array
      result[key] = Array.isArray(result[key]) ? [...result[key], parseValue(value)] : [result[key], parseValue(value)];
    } else {
      result[key] = parseValue(value);
    }
  }
  return result;
}

export function parseValue(value) {
  if (value === "true") return true;
  if (value === "false") return false;
  if (!isNaN(value) && value.trim() !== "") return Number(value);
  return value;
}

export function createSearchParameters(params) {
  const searchParams = new URLSearchParams();
  for (const key in params) {
    if (Array.isArray(params[key])) {
      params[key].forEach((value) => searchParams.append(key, value));
    } else {
      searchParams.set(key, params[key]);
    }
  }
  return searchParams.toString();
}
