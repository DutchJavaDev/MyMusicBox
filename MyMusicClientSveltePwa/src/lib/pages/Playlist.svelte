<script>
  // @ts-nocheck
  import { writable } from "svelte/store";
  import { onDestroy, onMount, setContext } from "svelte";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { playOrPauseSong, setPlaylists, updateCurrentPlaylist } from "../scripts/playbackService";
  import SongComponent from "../components/SongComponent.svelte";

  const updateIntervalTimeOut = 1000; // Update every second
  let intervalId

  export let playlistId = -1;

  let songs = writable([]);

  onMount(() => {
    songs.set(getCachedPlaylistSongs(playlistId));
    setContext("playOrPauseSong", playOrPause);

    intervalId = setInterval(() => {
      songs.set(getCachedPlaylistSongs(playlistId));
      updateCurrentPlaylist(playlistId);
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
  <div class="playlist-page">
    <div class="playlist-songs">
      {#each $songs as song}
        <SongComponent {song} {playlistId} />
      {/each}
    </div>
  </div>
{:else}
  <p>No songs in playlist.</p>
{/if}
