<script>
// @ts-nocheck

  import { isPlaying, currentSong, 
           playPercentage, setCurrentTime, 
           nextSong, previousSong, 
           isShuffledEnabled, isLoopingEnabled, 
           toggleShuffle, playOrPauseSong,
           toggleLoop, isLoading } from "../scripts/playbackService";
  import { getImageUrl } from "../scripts/api";
  import { get } from "svelte/store";
  import { isTimerEnabled, timeLeft, toggleSleepTimer } from "../scripts/sleeptimerService";

  $: $currentSong;
  $: $isPlaying;
  $: $playPercentage;
  $: $isShuffledEnabled;
  $: $isLoopingEnabled;
  $: $isTimerEnabled;
  $: $timeLeft;
  $: $isLoading
  
  function togglePlay() {
    playOrPauseSong(get(currentSong).id);
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
{#if $currentSong && $currentSong.id !== -999} <!-- Ensure currentSong is valid -->
  <div class="modal fade" id="songControlModal" tabindex="-1" aria-labelledby="songControlModalLabel" aria-hidden="false">
    <div class="modal-dialog modal-fullscreen-sm-down">
      <div class="modal-content">
        <div class="modal-body">
          <div>
            <div class="row">
              <div class="col-12">
                <div class="container-fluid" style="height: 7rem;">
                    <p class="text-white" id="songControlModalLabel">{$currentSong.name}</p>
                </div>
              </div>
              <div class="col-12 text-center">
                <img loading="lazy" class="img-fluid border border-1 rounded rounded-2 mt-1" src={getImageUrl($currentSong.thumbnail_path)} alt="404" />
              </div>
              <div class="col-12">
                <input type="range" on:change={seekEvent} class="form-range mt-5" value={$playPercentage} min="0" max="100" step="1" />
              </div>
              <div class="col-12">
                <div class="row mt-4">
                  <div class="col-4">
                    <button aria-label="previous song" on:click={previousSong} class="btn w-100">
                      <i class="fa-solid fa-backward fa-2xl"></i>
                    </button>
                  </div>

                  <div class="col-4">
                    <button on:click={togglePlay} class="btn w-100">
                      {#if $isPlaying}
                        <i class="fa-solid fa-pause fa-2xl"></i>
                      {:else if $isLoading}
                        <i class="fa-solid fa-spinner fa-spin fa-2xl"></i>
                      {:else}
                        <i class="fa-solid fa-play fa-2xl"></i>
                      {/if}
                    </button>
                  </div>

                  <div class="col-4">
                    <button aria-label="next song" on:click={nextSong} class="btn w-100">
                      <i class="fa-solid fa-forward fa-2xl"></i>
                    </button>
                  </div>
                </div>
              </div>
              <!-- Timer, shuffle and repeat controls -->
              <div class="col-12">
                <div class="row mt-5">
                  <div class="col-4">
                    <button on:click={toggleSleepTimer} aria-label="sleep timer" type="button" class="btn w-100">
                      <i class="fa-solid fa-stopwatch-20" style="{$isTimerEnabled ? "color: #1CC558;" : "color:white;"}">
                        <span style="font-size: 0.8rem;">
                            &nbsp;{$isTimerEnabled ? $timeLeft : ""}
                      </span>
                    </i>
                    </button>
                  </div>

                  <div class="col-4">
                    <button on:click={toggleShuffle} aria-label="shuffle playlist" type="button" class="btn w-100">
                      <i class="fa-solid fa-shuffle" style="{$isShuffledEnabled ? "color: #1CC558;" : "color:white;"}"></i>
                    </button>
                  </div>
                  <div class="col-4">
                    <button on:click={toggleLoop} aria-label="repeat song" type="button" class="btn w-100">
                      <i class="fa-solid fa-repeat" style="{$isLoopingEnabled ? "color: #1CC558;" : "color:white;"}"></i>
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
    border-color: #1CC558 !important;
  }

  p{
    font-size: 1rem !important;
    font-weight: bolder;
    color: white;
    text-align: center;
  }

  i{
    color: #1CC558;
    font-weight: bolder;
  }

  .modal-footer{
    border: none !important;
  }

  .modal-footer button {
    border-color: #1CC558 !important;
    background-color: #343a4000 !important;
  }


input[type="range"]::-webkit-slider-thumb {
   background-color: #1CC558;
}

input[type="range"]::-webkit-slider-runnable-track {
   background-color: #ACACAC;
}

.modal-content {
    background-color: #121212 !important;
    color: white;
  }

  .modal-body {
    padding: 1rem;
  }

  .form-range {
    color: #1CC558;
  }
</style>
