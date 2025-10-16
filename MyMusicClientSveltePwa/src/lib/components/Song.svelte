<script>
  // @ts-nocheck

  import { getContext, onMount, setContext } from "svelte";
  import { currentSong, isPlaying, isLoading } from "../scripts/playbackService";
  import { getImageUrl, deleteSongFromPlaylist } from "../scripts/api";
  import { removeSongFromPlaylist } from "../scripts/playlistService";

  export let song;
  export let playlistId;

  let playOrPauseSong;

  onMount(() => {
    playOrPauseSong = getContext("playOrPauseSong");
  });

  async function deleteSong() {
    if(confirm(`Are you sure you want to delete the song "${song.name}" from this playlist?`)) {
      var deleted = await deleteSongFromPlaylist(playlistId, song.id);

      if(deleted) {
        removeSongFromPlaylist(song.id);
      } else {
        alert(`Failed to delete song with ID: ${song.id} from playlist ID: ${playlistId}`);
      }
    }
  }

  $: $isPlaying;
  $: $currentSong;
  $: $isLoading;
</script>

{#if song}
  <div class="song" style="--url: url({getImageUrl(song.thumbnail_path)});">
    <div class="blur"></div>
    <div class="row align-items-center mt-3 content">
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
  </div>
{:else}
  <p>No song available.</p>
{/if}

<style>
  .song {
    position: relative;
    background: #2c2c2c;
    border-radius: 10px;
    padding: 0px 10px;
    margin-bottom: 10px;
    min-height: 80px;
    box-shadow: 0 4px 12px rgba(90, 89, 89, 0.4);
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    overflow: hidden;
    transition:
      transform 0.2s ease,
  }

  .song:hover {
    transform: translateY(-3px);
  }

  .song-info {
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
  }
  .title {
    color: #ffffff;
    font-weight: bold;
    font-size: 0.5rem;
    margin: 5px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .blur {
    position: absolute;
    inset: 0;
    background-color: rgba(0, 0, 0, 0.60);
    backdrop-filter: blur(5px);
    z-index: 1;
    border-radius: 5px;
  }

  /* Keep content visible above blur */
  .content {
    position: relative;
    z-index: 2;
  }
</style>
