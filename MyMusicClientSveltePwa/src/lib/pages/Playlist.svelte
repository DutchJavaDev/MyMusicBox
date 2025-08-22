<script>
  // @ts-nocheck
  import { writable } from "svelte/store";
  import { onDestroy, onMount, setContext } from "svelte";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { playOrPauseSong, setPlaylists, updateCurrentPlaylist } from "../scripts/playbackService";
  import SongComponent from "../components/SongComponent.svelte";

  const updateIntervalTimeOut = 1500; // Update every second
  let intervalId
  let updating = false

  export let playlistId = -1;

  let songs = writable([]);

  onMount(() => {
    songs.set(getCachedPlaylistSongs(playlistId));
    setContext("playOrPauseSong", playOrPause);

    intervalId = setInterval(() => {

      if (updating) return; // Prevent multiple updates at the same time
      updating = true;
      songs.set(getCachedPlaylistSongs(playlistId));
      updateCurrentPlaylist(playlistId);
      updating = false;
    }, updateIntervalTimeOut);
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
      {#each $songs as song}
        <div class="col-12 col-lg-4">
          <SongComponent {song} {playlistId} />
        </div>
      {/each}
    </div>
{:else}
  <p>No songs in playlist.</p>
{/if}
