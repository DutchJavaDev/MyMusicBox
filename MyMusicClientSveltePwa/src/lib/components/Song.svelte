<script>
  // @ts-nocheck

  import { getContext, onMount, setContext } from "svelte";
  import { currentSong, isPlaying, isLoading } from "../scripts/playbackService";
  import { getImageUrl } from "../scripts/api";

  export let song;

  let playOrPauseSong;

  onMount(() => {
    playOrPauseSong = getContext("playOrPauseSong");
  });

  $: $isPlaying;
  $: $currentSong;
  $: $isLoading;
</script>

{#if song}
  <div class="song">
  <div class="thumb" style="--url: url({getImageUrl(song.thumbnail_path)});"></div>
  <div class="song-info">
    <div class="title">{song.name}</div>
  </div>
  <button on:click={() => playOrPauseSong(song.id)} class="play-btn" title="Play">
    {#if $currentSong && $currentSong.id === song.id && $isPlaying}
          <i class="fa-solid fa-pause"></i>
        {:else if $isLoading && $currentSong.id === song.id}
          <i class="fa-solid fa-spinner fa-spin"></i>
        {:else}
          <i class="fa-solid fa-play"></i>
        {/if}
  </button>
</div>
{:else}
  <p>No song available.</p>
{/if}

<style>
  .song {
    display:flex;
    align-items:center;
    background:#2C2C2C;
    border-radius:12px;
    padding:10px 14px;
    margin-bottom:12px;
    box-shadow:0 4px 12px rgba(0,0,0,0.4);
    transition:transform 0.2s ease, background-color 0.2s ease;
  }
  .song:hover {
    background:#333333;
    transform:translateY(-2px);
  }
  .thumb {
    width:48px;
    height:48px;
    border-radius:8px;
    background-image: var(--url);
    background-size: cover;
    background-position: center;
    flex-shrink:0;
  }
  .song-info {
    flex:1;
    margin:0 12px;
    display:flex;
    flex-direction:column;
  }
  .title {
    color:#FFFFFF;
    font-weight:600;
    font-size: 0.8rem;
    margin: 5px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .play-btn {
    /* background:#1DB954; */
    border:none;
    border-color: transparent;
    width:36px;
    height:36px;
    color:#1DB954;
    background-color: transparent;;
    font-size:1.3rem;
    cursor:pointer;
    /* transition: background-color 0.2s; */
  }
  /* .play-btn:hover {
    background:#17a84b;
  } */
</style>
