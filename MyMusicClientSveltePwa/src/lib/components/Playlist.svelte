<script>
  // @ts-nocheck
  import { getImageUrl } from "../scripts/api";
  import { navigateTo } from "../scripts/routeService.js";
  import { playOrPauseSong, setPlaylists } from "../scripts/playbackService";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { onMount } from "svelte";
  // @ts-nocheck
  export let playlist = null;

  let songCount = 0;

  function viewPlaylist() {
    navigateTo(`/Playlists`, { playlistId: playlist.id });
  }

  function playPlaylist() {
    setPlaylists(playlist.id);
    playOrPauseSong(null);
  }

  onMount(() => {
    songCount = getCachedPlaylistSongs(playlist.id).length;
  });
</script>
{#if playlist}
    <!-- Playlist Card -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div on:click={viewPlaylist} aria-labelledby="view playlist" class="cursor">
      <article class="card" aria-labelledby={playlist.name}>
        <div class="art" style="--url: url({getImageUrl(playlist.thumbnailPath)});"></div>
        <div class="meta">
          <h3 id="playlist1">{playlist.name}</h3>
          <p>{songCount} songs • {playlist.description}</p>
        </div>
      </article>
    </div>
  {:else}
    <p>No playlist founnd with id: {playlist.id}.</p>
  {/if}

<style>
  .card {
    background: #2c2c2c;
    border-radius: 14px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 180px;
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.55);
    border: 1px solid rgba(0, 0, 0, 0.35);
    transition:transform 0.2s ease, background-color 0.2s ease;
  }

  .card .art {
    height: 120px;
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    
  }
  
  .card:hover {
    transform:scale(1.05);
  }

  .card .meta {
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .card h3 {
    margin: 0;
    font-size: 20px;
    color: #ffffff;
  }

  .card p {
    margin: 0;
    color: #b3b3b3;
    font-size: 14px;
  }
</style>
