<!-- App.svelte -->
<script>
  import Home from "./lib/pages/Home.svelte";
  import { onMount } from "svelte";
  import { initializeRoute, route, setRoute } from "./lib/route.js";
  import Test from "./lib/pages/Test.svelte";

  $: $route;

  onMount(() => {
    initializeRoute();
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
      {#if $route === "Home"}
        <Home />
      {:else if $route === "Test"}
        <Test />
      {:else}
        <div class="content-item">Page not found {$route}</div>
      {/if}
    </div>
  </main>

  <!-- Sticky Player Bar -->
  <div class="player-bar bg-dark rounded rounded-3 mb-2">
    <div class="row">
      <div class="col-3">
        <img class="img-fluid" src="https://www.bamdevserver.nl/dev/api/v1/images/_BJ9AJ7iwIg.jpg" alt="lol" />
      </div>
      <div class="col-7 bg-dark border border-1">
        <button class="btn btn-dark btn-lg clickable-text">Lorem ipsum dolor si Lorem ipsum dolor si Lorem ipsum dolor si Lorem ipsum dolor  Lorem ipsum dolor si Lorem ipsum dolor si</button>
      </div>
      <div class="col-2">
        <button class="btn btn-dark play-button">▶</button>
      </div>
    </div>
  </div>

  <!-- Sticky Bottom Bar -->
  <footer class="bottom-bar bg-dark">
    <div class="row w-100">
      <button class="col-4 btn btn-dark" on:click={() => setRoute("Settings")}>Settings</button>
      <button class="col-4 btn btn-dark" on:click={() => setRoute("Home")}>Home</button>
      <button class="col-4 btn btn-dark" on:click={() => setRoute("Test")}>Test</button>
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

  .player-bar {
    max-height: 10rem;
  }

  .player-bar .clickable-text {
    font-size: 1.1rem;
    max-height: 3.6rem;
    font-weight: bolder;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: start;
  }

  .play-button {
    font-weight: bolder;
    font-size: 2.5rem;
    width: 100%;
    height: 100%;
    color: #5bbd99;
    border-radius: 0.25rem;
    border: none;
    font-weight: bolder;
  }

  .bottom-bar button {
    font-size: 1.2rem;
    max-height: 3rem;
    border: none !important;
  }

  .player-bar img {
    width: 4.5rem;
    height: 4.5rem;
    object-fit: contain;
  }

  .top-bar {
    flex: 0 0 auto;
    padding: 1rem;
    position: sticky;
    height: 150px;
    top: 0;
    z-index: 10;
    text-align: center;
    border-bottom: 10px solid #ffffff;
    border-bottom-left-radius: 25px;
    border-bottom-right-radius: 25px;
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
    border-top: 15px solid #ffffff;
    border-top-left-radius: 25px;
    border-top-right-radius: 25px;
    height: 10rem; /* Optional: define fixed height if needed for padding calc */
  }

  @media (min-width: 600px) {
    .content-item {
      font-size: 1.1rem;
    }
  }
</style>
