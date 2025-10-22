<script>
  // @ts-nocheck
  import { get, writable } from "svelte/store";
  import { onDestroy, onMount, setContext } from "svelte";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { playOrPauseSong, setPlaylists, updateCurrentPlaylist } from "../scripts/playbackService";
  import SongComponent from "../components/Song.svelte";
  import VirtualList from "@sveltejs/svelte-virtual-list/VirtualList.svelte";
  import Song from "../components/Song.svelte";

  const updateIntervalTimeOutInMs = 750; // Update every 750 ms
  let intervalId;
  let updating = false;
  let start;
  let end;

  export let playlistId = -1;
  let visibleCount = 100;

  function loadMore() {
    visibleCount += 100;
  }

  let songs = writable([]);
  let readableSongs = [];

  onMount(() => {
    songs.set(getCachedPlaylistSongs(playlistId));
    setContext("playOrPauseSong", playOrPause);
    readableSongs = getCachedPlaylistSongs(playlistId);

    intervalId = setInterval(() => {
      // if (updating) return; // Prevent multiple updates at the same time
      updating = true;
      readableSongs = getCachedPlaylistSongs(playlistId);
      songs.set(readableSongs);
      updateCurrentPlaylist(playlistId);
      updating = false;
    }, updateIntervalTimeOutInMs);
  });

  function playOrPause(songId) {
    setPlaylists(playlistId);
  }

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

{#if readableSongs.length > 0}
<!-- <p>showing items {start}-{end}</p> -->
<div class='container'>
	<VirtualList items={readableSongs} bind:start bind:end let:item>
    <SongComponent song={item} {playlistId} />
	</VirtualList>
</div>
{:else}
  <p>No songs in playlist.</p>
{/if}

<style>
.container {
		border-top: 1px solid #333;
		border-bottom: 1px solid #333;
		min-height: 65vh;
		height: calc(100vh - 15em);
	}
</style>
