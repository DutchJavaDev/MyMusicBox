<script>
import { playAudio } from "../scripts/playback";
export let song;

function playSong() {
  if (song) {
    playAudio(song.source_id, song.id);
  }
}
</script>

{#if song}
  <!-- <div class="song-component bg-dark d-flex align-items-center mb-2 border border-1 rounded rounded-2 p-2">
        <div class="song-info">
            <p>{song.name}</p>
        </div>
        <button class="btn btn-dark play-button">▶</button>
    </div> -->
  <div class="row mb-3 song-component">
    <div class="col-10 bg-dark border border-1 rounded rounded-2">
      <div class="text-lg-start">
        <p>{song.name}</p>
      </div>
    </div>
    <div class="col-2">
      <button on:click={playSong} class="btn btn-dark play-button">▶</button>
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
