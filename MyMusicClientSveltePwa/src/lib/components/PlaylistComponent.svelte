<script>
  // @ts-nocheck
  export let playlist = null;
  import { getImageUrl, getPlaylistSongs } from "../scripts/api";
  import { setRoute } from "../scripts/route.js";
  import { currentPlaylistId } from "../scripts/playlist.js";

  $: $currentPlaylistId;

  async function viewPlaylist() {
    setRoute(`/Playlist`, { playlistId : playlist.id });
  }
</script>

<div class="playlist-component">
  {#if playlist}
    <button on:click={viewPlaylist} class="playlist-item btn w-100 border border-3" style="--url: url({getImageUrl(playlist.thumbnailPath)}); {$currentPlaylistId && $currentPlaylistId === playlist.id ? "border-color: #5bbd99 !important;" : ""}">
      <h3>{playlist.name}</h3>
      <p>{playlist.description}</p>
    </button>
  {:else}
    <p>No playlist available.</p>
  {/if}
</div>

<style>
  .playlist-item {
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    font-weight: bolder;
    color: white;
    height: 10rem;
  }

  .playlist-component {
    margin-top: 10px;
  }
</style>
