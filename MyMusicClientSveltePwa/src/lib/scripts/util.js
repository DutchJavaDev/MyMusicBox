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
  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
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
