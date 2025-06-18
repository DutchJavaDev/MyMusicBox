<script>
  // @ts-nocheck

  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import SongComponent from "../components/SongComponent.svelte";
  import { getPlaylistSongs } from "../scripts/api";

  let songs = writable([]);
  export let id = -1;

  $: $songs;

  onMount(() => {
    songs.set(getPlaylistSongs(id));
  });
</script>

{#if $songs.length > 0}
  <div class="playlist-page">
    <!-- <h1>{playlist.name}</h1>
    <p>{playlist.description}</p> -->
    <div class="playlist-songs">
      {#each $songs as song}
        <SongComponent {song} />
      {/each}
    </div>
  </div>
{:else}
  <p>No playlist selected.</p>
{/if}
