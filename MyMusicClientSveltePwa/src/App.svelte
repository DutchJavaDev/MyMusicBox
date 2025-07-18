<!-- App.svelte -->
<script>
  import { onMount, onDestroy } from "svelte";
  import { initializeRouteService, pathName, navigateTo, component, componentParams } from "./lib/scripts/routeService.js";
  import PlayerBarComponent from "./lib/components/PlayerBarComponent.svelte";
  import Modals from "./lib/components/Modals.svelte";
  import { initializePlaylistService } from "./lib/scripts/playlistService.js";
  import { initializePlaybackService } from "./lib/scripts/playbackService.js";
  import { initializeMediaSessionService } from "./lib/scripts/mediasessionService.js";

  $: $pathName;
  $: $component;

  // @ts-ignore
  const version = __APP_VERSION__;

  onMount(() => {
    async function initializeServices() {
      initializeRouteService();
      await initializePlaylistService();
      initializePlaybackService();
      initializeMediaSessionService();
    }
    initializeServices();
  });

  // This is a temporary function to handle refresh logic.
  // It can be replaced with a more specific implementation later.
  async function refresh() {}
</script>

<div class="app-layout bg-dark">
  <!-- Sticky Top Bar -->
  <header class="top-bar">
    <div class="container-fluid h-100">{$pathName} <span style="font-size: 0.8rem;">(v{version})</span></div>
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
        <button aria-label="empty storage" class="btn btn-dark w-100" on:click={refresh}><i class="fa-solid fa-trash"></i></button>
      </div>
      <div class="col-6">
        <button aria-label="home" class="btn btn-dark w-100" on:click={() => navigateTo("/Home")}><i class="fa-solid fa-house"></i></button>
      </div>
    </div>
  </footer>
</div>

<Modals />

<audio id="audio-player" preload="none" style="display: none;"></audio>

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
    height: 3.8rem; /* Optional: define fixed height if needed for padding calc */
  }

  .bottom-bar button {
    font-weight: bolder;
    border: 0.1rem solid #5bbd99 !important;
    background-color: #343a40 !important;
  }
</style>
