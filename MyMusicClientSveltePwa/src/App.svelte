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
  async function refresh() {
    window.location.reload();
  }
</script>

<div class="app-layout">
  <!-- Sticky Top Bar -->
  <header class="top-bar">
    <div class="top-bar-title text-center">MyMusicBox<span style="font-size: 0.8rem;">(v{version})</span></div>
    <div class="row">
      <div class="col-12 mt-2">
        <!-- Search Bar -->
        <div class="search" role="search" aria-label="Search music">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <path stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35" />
            <circle cx="11" cy="11" r="6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <input type="search" placeholder="Search music #not working jet..." aria-label="Search" />
        </div>
      </div>
    </div>
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
    <div class="row w-100 justify-content-center">
      <div class="col-3 col-lg-2 col-md-2 col-sm-2">
        <button aria-label="empty storage" class="btn btn-dark w-100 text-center" on:click={refresh}><i class="fa-solid fa-arrows-rotate"></i></button>
      </div>
      <div class="col-3 col-lg-2 col-md-2 col-sm-2">
        <button aria-label="home" class="btn btn-dark w-100" on:click={() => navigateTo("/Home")}><i class="fa-solid fa-house"></i></button>
      </div>
      <div class="col-3 col-lg-2 col-md-2 col-sm-2">
        <button aria-label="home" class="btn btn-dark w-100" on:click={() => navigateTo("/Settings")}><i class="fa-solid fa-gear"></i></button>
      </div>
    </div>
  </footer>
</div>

<Modals />

<audio id="audio-player" preload="none" style="display: none;"></audio>

<style>
  .search {
    background: #2a2a2a;
    padding: 14px 18px;
    border-radius: 28px;
    display: flex;
    align-items: center;
    gap: 12px;
    box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.4);
    width: 95%;
    margin: 0px auto;
  }
  .search svg {
    color: #b3b3b3;
    flex-shrink: 0;
  }
  .search input {
    background: transparent;
    border: 0;
    outline: 0;
    color: #ffffff;
    font-size: 16px;
    width: 100%;
  }
  .search input::placeholder {
    color: #b3b3b3;
  }

  .app-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    background-color: #1e1e1e;
  }

  .bottom-bar button {
    font-size: 0.6rem;
    max-height: 3rem;
    border: none !important;
  }
  .top-bar-title {
    font-size: 1.3rem;
    font-weight: bold;
    color:  #b3b3b3;
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
    border-top: 0.1rem solid #867878;
    background-color: unset;
    height: 3.2rem; /* Optional: define fixed height if needed for padding calc */
  }

  .bottom-bar button {
    font-weight: bolder;
    background-color: #2c2c2c !important;
  }
</style>
