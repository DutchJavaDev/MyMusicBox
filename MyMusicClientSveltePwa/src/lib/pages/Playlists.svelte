<script>
  // @ts-nocheck
  import { get, writable } from "svelte/store";
  import { getContext, onDestroy, onMount, setContext } from "svelte";
  import { getCachedPlaylistSongs } from "../scripts/storageService";
  import { playOrPauseSong, setPlaylists, updateCurrentPlaylist } from "../scripts/playbackService";
  import SongComponent from "../components/Song.svelte";
  import VirtualList from "@sveltejs/svelte-virtual-list/VirtualList.svelte";
  import Song from "../components/Song.svelte";
  import { searchQuery } from "../scripts/util.js";

  const updateIntervalTimeOutInMs = 250; // Update every 250 ms 1/4 of  a second
  let intervalId;
  let updating = false;
  let start;
  let end;

  let songs = writable([]);
  let readableSongs = [];
  let visibleSongs = writable([]);
  let searchQueryUnsubscribe;
  let filter = "";

  let componentParameters = $props();
  let playlistId = componentParameters["playlistId"];

  onMount(() => {
    searchQueryUnsubscribe = searchQuery.subscribe((value) => {
      if (!value) {
        filter = "";
        visibleSongs.set(readableSongs);
        return;
      }

      visibleSongs.set(readableSongs.filter((song) => song.name.toLowerCase().includes(value.toLowerCase())));

      filter = value;
    });

    let auto = getContext("preventAutoScroll");
    auto(null);

    songs.set(getCachedPlaylistSongs(playlistId));
    setContext("playOrPauseSong", playOrPause);
    readableSongs = getCachedPlaylistSongs(playlistId);

    // This is javascript and it might break, it might not until we look in to the box like shrowdinger's cat, but for bugs
    // @me-not-having-the-enegy-to-write-a-better-filter-function
    if (filter && filter.length > 0) {
      visibleSongs.set(readableSongs.filter((song) => song.name.toLowerCase().includes(filter.toLowerCase())));
    } else {
      visibleSongs.set(readableSongs);
    }

    intervalId = setInterval(() => {
      if (updating) return; // Prevent multiple updates at the same time
      updating = true;
      readableSongs = getCachedPlaylistSongs(playlistId);

      // This is javascript and it might break, it might not until we look in to the box like shrowdinger's cat, but for bugs
      // @me-not-having-the-enegy-to-write-a-better-filter-function
      if (filter && filter.length > 0) {
        visibleSongs.set(readableSongs.filter((song) => song.name.toLowerCase().includes(filter.toLowerCase())));
      } else {
        visibleSongs.set(readableSongs);
      }

      songs.set(readableSongs);
      updateCurrentPlaylist(playlistId);
      updating = false;
    }, updateIntervalTimeOutInMs);
  });

  function playOrPause(songId) {
    setPlaylists(playlistId);
    playOrPauseSong(songId);
  }

  onDestroy(() => {
    clearInterval(intervalId);
    if (searchQueryUnsubscribe) searchQueryUnsubscribe(); // stop listening
  });
</script>

{#if $visibleSongs && $visibleSongs.length > 0}
  <!-- <p>showing items {start}-{end}</p> -->
  <div class="container-cs border border-dark">
    <VirtualList items={$visibleSongs} let:item>
      <SongComponent song={item} playlistId={componentParameters["playlistId"]} />
    </VirtualList>
  </div>
{:else}
  <p>No songs found in playlist.</p>
{/if}

<style>
  .container-cs {
    background: #2a2a2a;
    border-top: 1px solid #333;
    border-bottom: 1px solid #333;
    /* min-height: calc(100vh - 30vh); */
    height: calc(100vh - 35vh);
    padding-right: unset;
    padding-left: unset;
  }
</style>
