<script>
  // @ts-nocheck

  import { getContext, onMount, setContext } from "svelte";
  import { currentSong, isPlaying, isLoading } from "../scripts/playbackService";
  import { getImageUrl, deleteSongFromPlaylist } from "../scripts/api";
  import { removeSongFromPlaylist } from "../scripts/playlistService";
  import { getCurrentPlaylistId, setCurrentSongIndex, setCurrentSongTime } from "../scripts/storageService";
  import { get } from "svelte/store";

  $: $isPlaying;
  $: $currentSong;
  $: $isLoading;

  export let song;
  export let playlistId;

  let playOrPauseSong;

  onMount(() => {
    playOrPauseSong = getContext("playOrPauseSong");
  });

  async function deleteSong() {
    if (confirm(`Are you sure you want to delete the song "${song.name}" from this playlist?`)) {
      var deleted = await deleteSongFromPlaylist(playlistId, song.id);

      if (deleted) {
        removeSongFromPlaylist(song.id);
        const currentPlaylistId = getCurrentPlaylistId();
        // If the deleted playlist is the current playing playlist, stop playback
        const _currentSong = get(currentSong);
        if (song.id === _currentSong.id) {
          // stop playback
          playOrPauseSong(null);
          currentSong.set({ id: -999, title: "", artist: "", album: "", source_id: "" });
          setCurrentSongTime(0);
          setCurrentSongIndex(-1);
        }
      } else {
        alert(`Failed to delete song with ID: ${song.id} from playlist ID: ${playlistId}`);
      }
    }
  }

</script>

{#if song}
  <!-- <div class="song" style="--url: url({getImageUrl(song.thumbnail_path)});">
    <div class="blur"></div>
    <div class="row align-items-center content">
      <div class="col-2">
        <button on:click={() => playOrPauseSong(song.id)} style="background-color: transparent; border: none; color: #1db954;">
          {#if $currentSong && $currentSong.id === song.id && $isPlaying}
            <i class="fa-solid fa-pause"></i>
          {:else if $isLoading && $currentSong.id === song.id}
            <i class="fa-solid fa-spinner fa-spin"></i>
          {:else}
            <i class="fa-solid fa-play"></i>
          {/if}
        </button>
      </div>
      <div class="song-info col-8">
        <div class="title">{song.name}</div>
      </div>
      <div class="col-2">
        <button on:click={deleteSong} class="text-center" aria-label="settings" style="background-color: transparent; border: none; font-size: 1rem;">
          <i class="fa-solid fa-trash text-danger"></i>
        </button>
      </div>
    </div>
  </div> -->

    <div class="row song m-1 rounded rounded-1" style="--url: url({getImageUrl(song.thumbnail_path)});">
    <div class="col-12 blur">
    </div>
    <div class="col-12">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="row content" aria-label="bsd" aria-roledescription="action" type="button" on:click={() => playOrPauseSong(song.id)} style="cursor: pointer; align-items: center; padding: 10px;">
      <div class="col-10">
        <div class="title {($currentSong && $currentSong.id === song.id && $isPlaying) ? "iamcute" : ""}">{song.name}</div>
      </div>
      <div class="col-2">
        <button on:click={deleteSong} class="text-center" aria-label="settings" style="background-color: transparent; border: none; font-size: 1rem;">
          <i class="fa-solid fa-trash text-danger"></i>
        </button>
      </div>
      </div>
    </div>
  </div>
{:else}
  <p>No song available.</p>
{/if}

<style>
  .song {
    position: relative;
    background: #2c2c2c;
    border-radius: 10px;
    margin-bottom: 10px;
    min-height: 80px;
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    overflow: hidden;
    transition: transform 0.2s ease;
    border-radius: 2px;
    margin: auto;
  }

  .song:hover {
    transform: translateY(-3px);
  }

  /* .song-info {
    flex: 1;
    margin: 0 12px;
    display: flex;
    flex-direction: column;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  } */
  .title {
    color:  #b3b3b3;
    font-size: 0.7rem;
    margin-top: 5px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .blur {
    position: absolute;
    inset: 0;
    background-color: #2c2c2c76;
    backdrop-filter: blur(10px);
    z-index: 1;
    /* border-radius: 5px; */
  }

  /* Keep content visible above blur */
  .content {
    position: relative;
    z-index: 2;
  }

@keyframes lights {
  0% {
    color: #1DB954;
    text-shadow:
      0 0 0.35em hsla(140, 80%, 50%, 0.25),
      0 0 0.1em hsla(140, 80%, 60%, 0.2);
  }
  
  30% { 
    color: #1db954b5;
    text-shadow:
      0 0 0.45em hsla(140, 80%, 50%, 0.35),
      0 0 0.15em hsla(140, 80%, 60%, 0.25),
      -0.15em -0.05em 0.15em hsla(60, 90%, 60%, 0.15),
      0.15em 0.05em 0.15em hsla(200, 100%, 60%, 0.2);
  }
  
  40% { 
    color: #1db9549d;
    text-shadow:
      0 0 0.4em hsla(140, 80%, 50%, 0.3),
      0 0 0.12em hsla(140, 80%, 80%, 0.25),
      -0.1em -0.05em 0.1em hsla(60, 90%, 60%, 0.12),
      0.1em 0.05em 0.1em hsla(200, 100%, 60%, 0.15);
  }
  
  70% {
    color: #1db95464;
    text-shadow:
      0 0 0.35em hsla(140, 80%, 50%, 0.25),
      0 0 0.1em hsla(140, 80%, 60%, 0.2),
      0.15em -0.05em 0.12em hsla(60, 90%, 60%, 0.12),
      -0.15em 0.05em 0.12em hsla(200, 100%, 60%, 0.18);
  }
  
  100% {
    color: #1DB954;
    text-shadow:
      0 0 0.35em hsla(140, 80%, 50%, 0.25),
      0 0 0.1em hsla(140, 80%, 60%, 0.2);
  }
}

/* body {
  margin: 0;
  font: 100% / 1.5 Raleway, sans-serif;
  color: hsl(230, 100%, 95%);
  background: linear-gradient(135deg, hsl(230, 40%, 12%), hsl(230, 20%, 7%));
  height: 100vh;
  display: flex;
} */

.iamcute {
  /* margin: auto;
  font-size: 3.5rem;
  font-weight: 300; */
  animation: lights 2s 500ms linear infinite;
}
</style>
