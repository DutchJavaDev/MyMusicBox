import { writable } from "svelte/store";

export let searchQuery = writable('');

// Mulberry32 PRNG
function mulberry32(seed) {
  return function () {
    let t = (seed += 0x6d2b79f5);
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 16)) >>> 0) / 4294967296;
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
  const rng = mulberry32(generateSeed());
  let i = array.length;
  while (i > 1) {
    const j = Math.floor(rng() * i--);
    [array[i], array[j]] = [array[j], array[i]];
  }
  return array;
}

export function getSearchParameters() {
  const result = {};
  for (const [key, value] of new URLSearchParams(window.location.search)) {
    const parsed = parseValue(value);
    if (key in result) {
      result[key] = Array.isArray(result[key]) ? [...result[key], parsed] : [result[key], parsed];
    } else {
      result[key] = parsed;
    }
  }
  return result;
}

/**
 * Parses a string value to its appropriate type.
 * @param {string} value
 * @returns {boolean|number|string|null|undefined}
 */
export function parseValue(value) {
  if (value === "true") return true;
  if (value === "false") return false;
  if (value == null) return value;
  const trimmed = value.trim();
  if (/^-?\d+(\.\d+)?$/.test(trimmed)) return Number(trimmed);
  return value;
}

/**
 * Creates a URL search string from an object of parameters.
 * @param {Object} params
 * @returns {string}
 */
export function createSearchParameters(params) {
  const searchParams = new URLSearchParams();
  for (const key in params) {
    const val = params[key];
    if (val == null) continue;
    if (Array.isArray(val)) {
      for (const v of val) searchParams.append(key, v);
    } else if (typeof val === "object") {
      searchParams.set(key, JSON.stringify(val));
    } else {
      searchParams.set(key, val);
    }
  }
  return searchParams.toString();
}
