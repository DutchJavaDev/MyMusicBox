<script>
  // @ts-nocheck
  import { getImageUrl } from "../scripts/api";
  import { navigateTo } from "../scripts/routeService.js";
  import { playOrPauseSong, setPlaylists } from "../scripts/playbackService";
  // @ts-nocheck
  export let playlist = null;

  function viewPlaylist() {
    navigateTo(`/Playlist`, { playlistId: playlist.id });
  }

  function playPlaylist() {
    setPlaylists(playlist.id);
    playOrPauseSong(null);
  }
</script>

<div class="playlist-component">
  {#if playlist}
    <div class="row me-1">
      <button aria-label="play button" on:click={playPlaylist} class="col-2 play-btn">
        <i class="fa-solid fa-play"></i>
      </button>
      <button aria-label="playlist button" on:click={viewPlaylist} class="playlist-item btn border border-3 col-10" style="--url: url({getImageUrl(playlist.thumbnailPath)});"> </button>
      <div class="text-start col-12 cursor">
        <p>#{playlist.name}<br />{playlist.description}</p>
      </div>
    </div>
  {:else}
    <p>No playlist available.</p>
  {/if}
</div>

<style>
  p {
    font-size: 0.7rem !important;
  }

  .playlist-item {
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    border-color: transparent !important;
    font-weight: bolder;
    color: white;
    min-height: 5rem;
    transform: scale(1);
    transition: transform 0.5s;
  }

  .playlist-item:hover {
    transform: scale(1.05);
    transition: transform 0.5s;
  }

  .play-btn {
    font-size: 1.8rem;
    background-color: rgba(131, 131, 131, 0.068) !important;
    border: none !important;
    color: #1cc558 !important;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    border-left: #1cc558 3px solid !important;
    transform: scale(1);
    transition: transform 0.5s;
  }

  .play-btn:hover {
    transform: scale(1.05);
    transition: transform 0.5s;
  }

  .playlist-component {
    margin-top: 10px;
  }

  .row{
    border: #1cc55711 2px solid;
    border-radius: 10px;
  }

  .row:hover{
    transition: background-color 0.4s;
    background-color: rgba(128, 128, 128, 0.096);
    border-radius: 15px;
    border-bottom: #1cc557 3px solid;
  }
</style>
