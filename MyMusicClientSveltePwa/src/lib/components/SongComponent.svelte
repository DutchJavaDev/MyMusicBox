<script>
import { playOrPauseAudio, currentSong, isPlaying } from "../scripts/playback.js";
import { setSongs, setCurrentSong, getCurrentSong } from "../scripts/playlist.js"

export let song;
export let playlistId;

$: $currentSong;
$: $isPlaying;

function playSong() {
  setSongs(playlistId);
  setCurrentSong(song);
  playOrPauseAudio(getCurrentSong());
}
</script>

{#if song}
  <div class="row mb-3 song-component">
    <div class="col-10 bg-dark border border-1 rounded rounded-2" style="{$currentSong && $currentSong.id === song.id && $isPlaying ? "border-color:#5bbd99 !important;" : ""}">
      <div class="text-lg-start">
        <p style="{$currentSong && $currentSong.id === song.id && $isPlaying ? "color:#5bbd99;" : ""}">{song.name}</p>
      </div>
    </div>
    <div class="col-2">
      <button on:click={playSong} class="btn btn-dark play-button">
        {#if $currentSong && $currentSong.id === song.id && $isPlaying}
          <i class="fa-solid fa-pause"></i>
        {:else}
          <i class="fa-solid fa-play"></i>
        {/if}
      </button>
    </div>
  </div>
{:else}
  <p>No song available.</p>
{/if}

<style>
  .song-component {
    height: 2.5rem;
    color: white;
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
    color: #5bbd99;
    align-items: center;
    justify-content: center;
  }
</style>
