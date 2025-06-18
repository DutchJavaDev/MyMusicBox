<script>
  // @ts-nocheck
  export let playlist = null;
  import { getImageUrl, getPlaylistSongs } from "../scripts/api";
  import { setRoute } from "../scripts/route.js";

  async function viewPlaylist() {
    setRoute(`/Playlist`, { id : playlist.id });
  }
</script>

<div class="playlist-component">
  {#if playlist}
    <button on:click={viewPlaylist} class="playlist-item btn w-100 border border-1" style="--url: url({getImageUrl(playlist.thumbnailPath)});">
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
