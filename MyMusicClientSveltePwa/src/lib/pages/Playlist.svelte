<script>
  // @ts-nocheck

  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import SongComponent from "../components/SongComponent.svelte";
  import { getPlaylistSongs } from "../scripts/api";

  let songs = writable([]);
  export let playlistId = -1;

  $: $songs;

  onMount(() => {
    songs.set(getPlaylistSongs(playlistId));
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
