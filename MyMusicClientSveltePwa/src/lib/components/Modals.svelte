<script>
  // @ts-nocheck

  import { isPlaying, currentSong, playPercentage, setCurrentTime, nextSong, previousSong, isShuffledEnabled, isLoopingEnabled, toggleShuffle, playOrPauseSong, toggleLoop, isLoading } from "../scripts/playbackService";
  import { getImageUrl, createPlaylist } from "../scripts/api";
  import { get } from "svelte/store";
  import { isTimerEnabled, timeLeft, toggleSleepTimer } from "../scripts/sleeptimerService";

  $: $currentSong;
  $: $isPlaying;
  $: $playPercentage;
  $: $isShuffledEnabled;
  $: $isLoopingEnabled;
  $: $isTimerEnabled;
  $: $timeLeft;
  $: $isLoading;

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

  async function handleCreatePlaylistSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    var response = await createPlaylist(formData);

    if (response.success){
      alert(`Playlist has been created successfully.`);
      
      // Close modal
      const modalElement = document.getElementById('createPlaylistModal');
      const modalInstance = bootstrap.Modal.getInstance(modalElement);
      
      // clear form
      modalInstance.hide();
      
      // clear form
      event.target.reset();
    } else {
      alert(`Failed to create playlist: ${response.data}`);
    }
  }
</script>

<!-- Play Modal -->
{#if $currentSong && $currentSong.id !== -999}
  <!-- Ensure currentSong is valid -->
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
                <img loading="lazy" class="img-fluid rounded rounded-2 mt-1" src={getImageUrl($currentSong.thumbnail_path)} alt="404" />
              </div>
              <div class="col-12">
                <input type="range" on:change={seekEvent} class="form-range mt-5" value={$playPercentage} min="0" max="100" step="0.3" />
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
                      <i class="fa-solid fa-stopwatch-20" style={$isTimerEnabled ? "color: #1CC558;" : "color:#ACACAC;"}>
                        <span style="font-size: 0.8rem;">
                          &nbsp;{$isTimerEnabled ? $timeLeft : ""}
                        </span>
                      </i>
                    </button>
                  </div>

                  <div class="col-4">
                    <button on:click={toggleShuffle} aria-label="shuffle playlist" type="button" class="btn w-100">
                      <i class="fa-solid fa-shuffle" style={$isShuffledEnabled ? "color: #1CC558;" : "color:#ACACAC;"}></i>
                    </button>
                  </div>
                  <div class="col-4">
                    <button on:click={toggleLoop} aria-label="repeat song" type="button" class="btn w-100">
                      <i class="fa-solid fa-repeat" style={$isLoopingEnabled ? "color: #1CC558;" : "color:#ACACAC;"}></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-dark w-100 text-white" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- Create Playlist Modal -->
<div class="modal fade" id="createPlaylistModal" tabindex="-1" aria-labelledby="createPlaylistModalLabel" aria-hidden="false">
  <div class="modal-dialog modal-fullscreen-sm-down">
    <div class="modal-content">
      <div class="modal-body">
        <form on:submit|preventDefault={handleCreatePlaylistSubmit} class="p-2 rounded rounded-2 tile-bg">
          <div class="mb-3">
            <label for="playlistName" class="form-label">Playlist Name</label>
            <input type="text" required class="form-control form-control-sm" id="playlistName" name="playlistName" placeholder="Name of newly created playlist" />
          </div>
          <div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
            <label for="backgroundImage" class="form-label">Playlist Image (leave blank for default)</label>
            <input type="file" id="backgroundImage" name="backgroundImage" class="form-file-input form-control-sm" />
          </div>
          <div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
            <label for="publicPlaylist" class="form-label">Public</label>
            <input type="checkbox" checked id="publicPlaylist" name="publicPlaylist" class="form-check-input" />
          </div>
          <div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
            <label for="playlistDescription" class="form-label">Description</label>
            <textarea id="playlistDescription" name="playlistDescription" rows="3" class="form-control form-control-sm"> </textarea>
          </div>
          <button type="submit" class="btn btn-primary">Create Playlist</button>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-dark w-100 text-white" data-bs-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>

<style>
  .tile-bg {
    background-color: #2c2c2c;
  }
  img {
    height: 10rem;
    object-fit: contain;
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.55);
    border: 1px solid rgba(0, 0, 0, 0.35);
  }

  p {
    font-size: 1rem !important;
    font-weight: bolder;
    color: white;
    text-align: center;
  }

  i {
    font-size: 1.2rem;
    color: #acacac;
    font-weight: bolder;
  }

  .modal-footer {
    border: none !important;
  }

  .modal-footer button {
    background-color: #2c2c2c !important;
  }

  input[type="range"]::-webkit-slider-thumb {
    background-color: #1db954;
  }

  input[type="range"]::-webkit-slider-runnable-track {
    background-color: #acacac;
  }

  .modal-content {
    background-color: #1e1e1e !important;
    color: white;
  }

  .modal-body {
    padding: 1rem;
  }

  .form-range {
    color: #1db954;
  }
</style>
