<!-- App.svelte -->
<script>
  import { onMount } from "svelte";
  import { initializeRoute, route, setRoute, component, componentParams } from "./lib/scripts/route.js";
  import { updateStores } from "./lib/scripts/api.js";
  import { initPlaybackAudio } from "./lib/scripts/playback.js";

  $: $route;
  $: $component;

  onMount(() => {
    async function async() {
      await updateStores();
      initPlaybackAudio();
      initializeRoute();
    }
    async();
  });
</script>

<div class="app-layout">
  <!-- Sticky Top Bar -->
  <header class="top-bar">
    <div class="container-fluid h-100">{$route}</div>
  </header>

  <!-- Scrollable Content -->
  <main class="scrollable-content">
    <div class="container-fluid">
      <svelte:component this={$component} {...$componentParams} />
    </div>
  </main>

  <!-- Sticky Player Bar -->
  <div class="container-fluid player-bar mb-2 bg-dark rounded rounded-2">
    <div class="row justify-content-center align-items-center">
      <div class="col-10">
        <button class="btn btn-dark clickable-text">
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
        </button>
      </div>
      <div class="col-2 border-start border-2 rounded-start rounded-1">
        <button class="btn btn-dark play-button">▶</button>
      </div>
    </div>

  </div>

  <!-- Sticky Bottom Bar -->
  <footer class="bottom-bar">
    <div class="row w-100">
      <button class="col-4 btn btn-dark" on:click={() => setRoute("/Settings")}>Settings</button>
      <button class="col-4 btn btn-dark" on:click={() => setRoute("/Home")}>Home</button>
      <button class="col-4 btn btn-dark" on:click={() => setRoute("/Test")}>Test</button>
    </div>
  </footer>
</div>

<audio id="audio-player" style="display: none;"></audio>

<style>
  .app-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
  }

  .player-bar .clickable-text {
    font-size: 0.6rem;
    max-height: 3.6rem;
    font-weight: bold;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: start;
  }

  .play-button {
    font-weight: bolder;
    font-size: 2rem;
    width: 100%;
    height: 100%;
    display: block !important;
    margin: 0;
    color: #5bbd99;
    border-radius: 0.25rem;
    border: none;
    font-weight: bolder;
  }

  .player-bar .col-2 {
    padding: 0 !important;
  }

  .player-bar .col-10 {
    padding: 0 !important;
  }

  .bottom-bar button {
    font-size: 1.2rem;
    max-height: 3rem;
    border: none !important;
  }

  /* .player-bar img {
    width: 4.5rem;
    height: 4.5rem;
    object-fit: contain;
  } */

  .top-bar {
    flex: 0 0 auto;
    padding: 1rem;
    position: sticky;
    height: 3.5rem;
    top: 0;
    z-index: 10;
    text-align: center;
    border-bottom: 0.2rem solid #ffffff;
    border-bottom-left-radius: 1.5rem;
    border-bottom-right-radius: 1.5rem;
  }

  .scrollable-content {
    flex: 1 1 auto;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    padding: 1rem 1rem 3rem; /* 👈 Important: bottom padding to make space for bottom bar */
  }

  .bottom-bar {
    flex: 0 0 auto;
    padding: 0.5rem;
    position: sticky;
    bottom: 0;
    z-index: 10;
    display: flex;
    justify-content: center;
    border-top: 0.2rem solid #ffffff;
    border-top-left-radius: 1.5rem;
    border-top-right-radius: 1.5rem;
    height: 9rem; /* Optional: define fixed height if needed for padding calc */
  }

  .bottom-bar button {
    background-color: black;
    font-weight: bolder;
    border: 0.1rem solid #ffffff !important;
  }
</style>
