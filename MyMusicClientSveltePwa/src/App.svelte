<!-- App.svelte -->
<script>
  import { onMount } from "svelte";
  import { initializeRoute, route, setRoute, component, componentParams } from "./lib/scripts/route.js";
  import { updateStores } from "./lib/scripts/api.js";
  import { initPlaybackAudio, playOrPauseAudio } from "./lib/scripts/playback.js";
  import { nextSong, previousSong } from "./lib/scripts/playlist.js";
  // @ts-ignore
  import {initializeMediaSession} from "./lib/scripts/mediasession.js";
  import PlayerBarComponent from "./lib/components/PlayerBarComponent.svelte";
  import Modals from "./lib/components/Modals.svelte";

 // @ts-ignore
  // @ts-ignore
    $: $route;
 // @ts-ignore
  // @ts-ignore
    $: $component;

  onMount(() => {
    async function async() {
      await updateStores();
      initPlaybackAudio();
      initializeMediaSession();
      initializeRoute();

      // setInterval( async () => {
      //   await updateStores();
      // }, 1000 * 30); // Update every 30 seconds
    }
    async();
  });

  // @ts-ignore
  // @ts-ignore
  function next(){
    playOrPauseAudio(nextSong())
  }
  // @ts-ignore
  // @ts-ignore
  function prev() {
    playOrPauseAudio(previousSong())
  }
</script>

<div class="app-layout bg-dark">
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
  <PlayerBarComponent />

  <!-- Sticky Bottom Bar -->
  <footer class="bottom-bar">
    <div class="row w-100">
      <div class="col-6">
        <button aria-label="settings" class="btn btn-dark w-100"><i class="fa-solid fa-gear"></i></button>
      </div>
      <div class="col-6">
        <button aria-label="home" class="btn btn-dark w-100" on:click={() => setRoute("/Home")}><i class="fa-solid fa-house"></i></button>
      </div>
    </div>
  </footer>
</div>

<Modals />

<style>
  .app-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
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
    border-bottom: 0.2rem solid #5bbd99;
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
    border-top: 0.2rem solid #5bbd99;
    border-top-left-radius: 1.5rem;
    border-top-right-radius: 1.5rem;
    height: 9rem; /* Optional: define fixed height if needed for padding calc */
  }

  .bottom-bar button {
    font-weight: bolder;
    border: 0.1rem solid #5bbd99 !important;
    background-color: #343a40 !important;
  }
</style>
