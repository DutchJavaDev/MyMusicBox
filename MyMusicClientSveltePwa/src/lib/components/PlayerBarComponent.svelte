<script>
  import { isPlaying, currentSong, playOrPauseAudio, playPercentage } from "../scripts/playback.js";

  $: $isPlaying;
  $: $currentSong;
  $: $playPercentage;

  function togglePlay() {
    playOrPauseAudio(null);
  }
</script>

<div class="container-fluid player-bar mb-2 rounded rounded-5">
  <div class="row space-between">
    <div class="col-9 rounded-end rounded-end-0 rounded-5 border border-1 border-white" style="background: linear-gradient(to right, gray {$playPercentage}%, #5bbd99 {$playPercentage}%);">
      <button type="button" class="btn clickable-text rounded-end rounded-end-0 rounded-5" data-bs-toggle="modal" data-bs-target="#songControlModal">
        {#if $currentSong}
          {$currentSong.name}
        {:else}
          No song playing
        {/if}
      </button>
    </div>
    <div class="col-3 border-start border-2">
      <button on:click={togglePlay} class="btn btn-dark border border-1 border-white play-button rounded-end rounded-end-5 w-100">
        {#if $isPlaying}
          <i class="fa-solid fa-pause"></i>
        {:else}
          <i class="fa-solid fa-play"></i>
        {/if}
      </button>
    </div>
  </div>
</div>

<style>
  .player-bar .clickable-text {
    font-size: 0.85rem;
    max-height: 2.8rem;
    min-height: 2.8rem;
    width: 100%;;
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

  .player-bar{
    background-color: gray !important;
  }

  .play-button {
    font-weight: bolder;
    font-size: 1.4rem;
    width: 100%;
    height: 100%;
    display: block !important;
    margin: 0;
    color: #5bbd99;
    font-weight: bolder;
  }

  .player-bar .col-9 {
    padding: 0 !important;
  }

  .player-bar .col-3 {
    padding: 0 !important;
  }
</style>
