<!-- App.svelte -->
<script>
// @ts-nocheck

  import { onMount, onDestroy, setContext } from "svelte";
  import { get, writable } from "svelte/store";
  import { initializeRouteService, pathName, navigateTo, component as currentView, componentParams } from "./lib/scripts/routeService.js";
  import PlayerBar from "./lib/components/PlayerBar.svelte";
  import Modals from "./lib/components/Modals.svelte";
  import { initializePlaylistService, deleteCurrentPlaylist } from "./lib/scripts/playlistService.js";
  import { initializePlaybackService } from "./lib/scripts/playbackService.js";
  import { initializeMediaSessionService } from "./lib/scripts/mediasessionService.js";
  import { searchQuery } from "./lib/scripts/util.js";
  import SearchBar from "./lib/components/SearchBar.svelte";

  $: $pathName;
  $: $currentView;

  // @ts-ignore
  export const version = __APP_VERSION__;

  let autoScroll = writable(true);

  export function preventAutoScroll(value = true) {
    // flips state of autoScroll, false to true and true to false and so on for eternity
    // that is all it should do.... but its javaScript so who knows what will happen in 2 weeks
    if (value && value === true || value === false) {
      autoScroll.set(value);
    }
    else{
      autoScroll.set(!get(autoScroll));
    }
  }

  onMount(() => {
    currentView.subscribe((value) => {
      // Whenever the component changes, enable auto scroll
     preventAutoScroll()
    });
    setContext("preventAutoScroll", preventAutoScroll);
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
      <!-- Configurable via js/storage/options etc per page / component -->
      <div class="col-12 mt-2 justify-content-center">
        <!-- Search Bar -->
         <SearchBar />
      </div>
    </div>
  </header>

  <!-- Scrollable Content -->
  <main class="{($autoScroll ? "scrollable-content" :"none-scrollable-content")}  ">
    <div class="container-fluid">
      <svelte:component this={$currentView} {...$componentParams} />
    </div>
  </main>

  <!-- Sticky Player Bar -->
  <PlayerBar />

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
        <div class="btn-group dropup w-100">
          
          <button type="button" aria-label="home" class="btn btn-dark w-100 dropdown-toggle" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false"><i class="fa-solid fa-plus"></i></button>
          
          <div id="actions" class="dropdown-menu w-100">
            <button class="btn btn-primary dropdown-item" data-bs-toggle="modal" data-bs-target="#createPlaylistModal" >New Playlist</button>
            <button class="btn btn-primary dropdown-item text-danger" on:click={deleteCurrentPlaylist} disabled={!$pathName.includes("playlists")}>Delete Current Playlist</button>
          </div>
        </div>
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
  #actions{
    background-color: #1e1e1e;
  }

  #actions button{
    padding: 5px auto;
    border: 5px solid red;
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

  .none-scrollable-content {
    flex: 1 1 auto;
    overflow-y: hidden;
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
    height: 3.2rem; 
  }
  /* Optional: define fixed height if needed for padding calc */

  .bottom-bar button {
    font-weight: bolder;
    background-color: #2c2c2c !important;
    color:  white;
  }
</style>
