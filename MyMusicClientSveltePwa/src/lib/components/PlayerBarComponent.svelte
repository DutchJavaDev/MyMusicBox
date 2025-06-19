<script>
  import { isPlaying, currentSong, playOrPauseAudio, playPercentage } from "../scripts/playback.js";

  $: $isPlaying;
  $: $currentSong;
  $: $playPercentage;

  function togglePlay() {
    playOrPauseAudio(null);
  }
</script>

{#if $isPlaying}
<div class="container-fluid player-bar mb-2 rounded rounded-5">
  <div class="row">
    <div class="col-10 rounded-end rounded-end-0 rounded-5" style="background: linear-gradient(to right, gray {$playPercentage}%, #5bbd99 {$playPercentage}%);">
      <button class="btn clickable-text">
        {#if $currentSong}
          {$currentSong.name}
        {:else}
          No song playing
        {/if}
      </button>
    </div>
    <div class="col-2 border-start border-2">
      <button on:click={togglePlay} class="btn play-button rounded-end rounded-end-5">{$isPlaying ? "||" : "▶"}</button>
    </div>
  </div>
</div>
{/if}
<style>
  .player-bar .clickable-text {
    font-size: 0.85rem;
    max-height: 2.8rem;
    font-weight: bold;
    color: white;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: center;
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

  .player-bar .col-2 {
    padding: 0 !important;
  }

  .player-bar .col-10 {
    padding: 0 !important;
  }
</style>
