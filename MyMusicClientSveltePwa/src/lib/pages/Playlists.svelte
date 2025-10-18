<script>
  // @ts-nocheck
  import { writable } from "svelte/store";
  import { onDestroy, onMount, setContext } from "svelte";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { playOrPauseSong, setPlaylists, updateCurrentPlaylist } from "../scripts/playbackService";
  import SongComponent from "../components/Song.svelte";

  const updateIntervalTimeOutInMs = 750; // Update every 750 ms
  let intervalId
  let updating = false

  export let playlistId = -1;
  let visibleCount = 100;

  function loadMore() {
    visibleCount += 100;
  }

  let songs = writable([]);

  onMount(() => {
    songs.set(getCachedPlaylistSongs(playlistId));
    setContext("playOrPauseSong", playOrPause);

    intervalId = setInterval(() => {
      // if (updating) return; // Prevent multiple updates at the same time
      updating = true;
      songs.set(getCachedPlaylistSongs(playlistId));
      updateCurrentPlaylist(playlistId);
      updating = false;
    }, updateIntervalTimeOutInMs);
  });

  function playOrPause(songId) {
    setPlaylists(playlistId);
    playOrPauseSong(songId);
  }

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

{#if $songs.length > 0}
  <div class="row">
      {#each $songs.slice(0, visibleCount) as song}
        <div class="col-12 col-lg-4">
          <SongComponent {song} {playlistId} />
        </div>
      {/each}
      {#if visibleCount < $songs.length}
        <div class="col-12 text-center my-3">
          <button class="btn btn-primary" on:click={loadMore}>Load More</button>
        </div>
      {/if}
    </div>
{:else}
  <p>No songs in playlist.</p>
{/if}
