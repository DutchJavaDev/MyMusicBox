<script>
  // @ts-nocheck

  import { onDestroy, onMount } from "svelte";
  import { writable } from "svelte/store";
  import SongComponent from "../components/SongComponent.svelte";
  import { getPlaylistSongs, writablePlaylistsStore } from "../scripts/api";
  import { setSongs } from "../scripts/playlist";
  import { on } from "svelte/events";

  let songs = writable([]);
  export let playlistId = -1;

  let intervalId;

  $: writablePlaylistsStore[playlistId];

  onMount(() => {
    songs.set(getPlaylistSongs(playlistId));
    writablePlaylistsStore[playlistId].subscribe((value) => {
      songs.set(value);
    });

  });

  onDestroy(() => {
    console.log("Cleaning up interval for playlist songs");
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
