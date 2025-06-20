<script>
  import { isPlaying, currentSong, playOrPauseAudio, playPercentage, toggleRepeat, isRepeatEnabled, setCurrentTime } from "../scripts/playback.js";
  import { getImageUrl } from "../scripts/api.js";
  import { nextSong, previousSong, shufflePlaylist, isShuffleEnabled } from "../scripts/playlist.js";

  $: $isPlaying;
  $: $currentSong;
  $: $playPercentage;
  $: $isShuffleEnabled;
  $: $isRepeatEnabled;

  function next() {
    playOrPauseAudio(nextSong());
  }
  function prev() {
    playOrPauseAudio(previousSong());
  }

  function togglePlay() {
    playOrPauseAudio(null);
  }

  function seekEvent(event) {
    const percentage = event.target.value;
    playPercentage.set(percentage);
    const currentSongData = $currentSong;
    if (currentSongData) {
      const duration = currentSongData.duration;
      const newTime = (duration * percentage) / 100;
      console.log(`Seeking to ${newTime} seconds in song ${currentSongData.name}`);
      setCurrentTime(newTime);
    }
  }
</script>

<!-- Modal -->
{#if $currentSong}
  <div class="modal fade" id="songControlModal" tabindex="-1" aria-labelledby="songControlModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-fullscreen-sm-down">
      <div class="modal-content bg-dark">
        <div class="modal-body">
          <div>
            <div class="row">
              <div class="col-12">
                <div class="container-fluid" style="height: 7rem;">
                    <p class="text-white" id="songControlModalLabel">{$currentSong.name}</p>
                </div>
              </div>
              <div class="col-12 text-center">
                <img class="img-fluid border border-1 rounded rounded-2 mt-1" src={getImageUrl($currentSong.thumbnail_path)} alt="404" />
              </div>
              <div class="col-12">
                <input type="range" on:change={seekEvent} class="form-range mt-5" value={$playPercentage} min="0" max="100" step="1" />
              </div>
              <div class="col-12">
                <div class="row mt-4">
                  <div class="col-4">
                    <button aria-label="previous song" on:click={prev} class="btn btn-dark w-100">
                      <i class="fa-solid fa-backward fa-2xl"></i>
                    </button>
                  </div>

                  <div class="col-4">
                    <button on:click={togglePlay} class="btn btn-dark w-100">
                      {#if $isPlaying}
                        <i class="fa-solid fa-pause fa-2xl"></i>
                      {:else}
                        <i class="fa-solid fa-play fa-2xl"></i>
                      {/if}
                    </button>
                  </div>

                  <div class="col-4">
                    <button aria-label="next song" on:click={next} class="btn btn-dark w-100">
                      <i class="fa-solid fa-forward fa-2xl"></i>
                    </button>
                  </div>
                </div>
              </div>
              <!-- Timer, shuffle and repeat controls -->
              <div class="col-12">
                <div class="row mt-5">
                  <div class="col-4">
                    <button disabled aria-label="sleep timer" type="button" class="btn btn-dark w-100">
                      <i class="fa-solid fa-stopwatch-20" style="color: white !important;">
                        <span style="font-size: 0.5rem;">
                            &nbsp;TODO
                      </span>
                    </i>
                    </button>
                  </div>

                  <div class="col-4">
                    <button on:click={shufflePlaylist} aria-label="shuffle playlist" type="button" class="btn btn-dark w-100">
                      <i class="fa-solid fa-shuffle" style="{$isShuffleEnabled ? "color: #5bbd99;" : "color:white;"}"></i>
                    </button>
                  </div>
                  <div class="col-4">
                    <button on:click={toggleRepeat} aria-label="repeat song" type="button" class="btn btn-dark w-100">
                      <i class="fa-solid fa-repeat" style="{$isRepeatEnabled ? "color: #5bbd99;" : "color:white;"}"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-dark fw-bolder w-100 text-white border border-2" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  img {
    height: 10rem;
    object-fit: contain;
    border-color: #5bbd99 !important;
  }

  p{
    font-size: 1rem !important;
    font-weight: bolder;
    color: white;
    text-align: center;
  }

  i{
    color: #5bbd99;
    font-weight: bolder;
  }

  .modal-footer{
    border: none !important;
  }

  .modal-footer button {
    border-color: #5bbd99 !important;
    background-color: #343a40 !important;
  }

input[type="range"]::-webkit-slider-thumb {
   background-color: #5bbd99;  
}
</style>
