<script>
  // @ts-nocheck
  import { get } from "svelte/store";
  import { isPlaying, currentSong, playOrPauseSong, playPercentage, isLoading } from "../scripts/playbackService";
  import { getImageUrl } from "../scripts/api";

  $: $currentSong;
  $: $isPlaying;
  $: $playPercentage;
  $: $isLoading;

  function togglePlay() {
    if(get(currentSong)){
      playOrPauseSong(get(currentSong).id);
    }
  }
</script>

<div class="container-fluid player-bar mb-2">
  <div class="row space-between">
    <div class="image-placeholder col-2 col-md-2 col-lg-2" style="--url: url({getImageUrl($currentSong.thumbnail_path)});">
      &nbsp;
    </div>
    <div class="col-8 col-md-8 col-lg-9" style="background: linear-gradient(to right, #1DB954 {$playPercentage}%, #2c2c2c {$playPercentage}%);">
      <button type="button" class="btn clickable-text" data-bs-toggle="{$currentSong ? "modal" : ""}" data-bs-target="{$currentSong ? "#songControlModal" : ""}">
        {#if $currentSong}
          {$currentSong.name}
        {:else}
          No song playing
        {/if}
      </button>
    </div>
    <div class="col-2 col-md-2 col-lg-1">
      <button on:click={togglePlay} class="btn play-button w-100">
        {#if $currentSong && $isPlaying && !$isLoading}
          <i class="fa-solid fa-pause"></i>
        {:else if !$isLoading && !$isPlaying}
          <i class="fa-solid fa-play"></i>
        {:else if $isLoading}
          <i class="fa-solid fa-spinner fa-spin"></i>
        {/if}
      </button>
    </div>
  </div>
</div>

<style>
  .player-bar .clickable-text {
    font-size: 0.65rem;
    max-height: 2.8rem;
    min-height: 2.8rem;
    width: 100%;
    font-weight: bold;
    color: white;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 2px;
  }

  .player-bar {
    background-color: transparent;
  }

  .image-placeholder{
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
  }

  .play-button {
    font-weight: bolder;
    font-size: 1.4rem;
    width: 100%;
    height: 100%;
    display: block !important;
    margin: 0;
    color: #1DB954;
    font-weight: bolder;
    background-color: #2c2c2c;
    border: none !important;
    border-radius: 0 !important;
  }

  .player-bar .col-8 {
    padding: 0 !important;
  }

  .player-bar .col-2 {
    padding: 0 !important;
  }
</style>
