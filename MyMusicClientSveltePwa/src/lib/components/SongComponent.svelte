<script>
  // @ts-nocheck

  import { getContext, onMount, setContext } from "svelte";
  import { currentSong, isPlaying, isLoading } from "../scripts/playbackService";

  export let song;

  let playOrPauseSong;

  onMount(() => {
    playOrPauseSong = getContext("playOrPauseSong");
  });

  $: $isPlaying;
  $: $currentSong;
  $: $isLoading;
</script>

{#if song}
  <div class="row mb-3 song-component">
    <div class="col-2">
      <button on:click={() => playOrPauseSong(song.id)} class="btn btn-dark w-100 play-button">
        {#if $currentSong && $currentSong.id === song.id && $isPlaying}
          <i class="fa-solid fa-pause"></i>
        {:else if $isLoading && $currentSong.id === song.id}
          <i class="fa-solid fa-spinner fa-spin"></i>
        {:else}
          <i class="fa-solid fa-play"></i>
        {/if}
      </button>
    </div>
    <div class="col-10 border border-1 rounded rounded-2 cursor" style={$currentSong && $currentSong.id === song.id ? "border-color:#1CC558 !important;" : "border-color: gray !important;"}>
      <div class="text-lg-start">
        <p style={$currentSong && $currentSong.id === song.id ? "color:#1CC558;" : ""}>{song.name}</p>
      </div>
    </div>
  </div>
{:else}
  <p>No song available.</p>
{/if}

<style>
  .song-component {
    height: 2.5rem;
    color: #B3B3B3;
    font-weight: bolder;
  }

  .song-component p {
    font-size: 0.8rem;
    margin: 5px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .play-button {
    font-size: 1.5rem;
    width: 3rem;
    height: 3rem;
    display: flex;
    color: #1cc558;
    background-color: transparent !important;
    border: none !important;
    align-items: center;
    justify-content: center;
    border-left: #1CC558 2px solid !important;
  }

  .play-button:hover {
    background-color: rgba(28, 197, 88, 0.1) !important;
    transition: background-color 0.3s ease;
  }
</style>
