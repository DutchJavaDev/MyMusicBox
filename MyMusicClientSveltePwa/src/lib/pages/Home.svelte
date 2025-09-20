<script>
  import { onDestroy, onMount } from "svelte";
  import PlaylistComponent from "../components/Playlist.svelte";
  import { playlistsStore } from "../scripts/playlistService.js";
  // @ts-ignore
  import { searchQuery } from "../scripts/util.js";

  // @ts-ignore
 // @ts-ignore
  //  $: $playlistsStore;

   let searchQueryUnsubscribe;
   let playlistsStoreUnsubscribe;

   let originalPlaylists = [];
   let filteredPlaylists = $state([]);

  onMount(() => {

    playlistsStoreUnsubscribe = playlistsStore.subscribe((value) => {
      originalPlaylists = value;
      filteredPlaylists = value;
    });

    // @ts-ignore
    searchQueryUnsubscribe = searchQuery.subscribe((value) => {
      if (value && value.length > 0){
        // filter playlists
        const lowerValue = value.toLowerCase();
        filteredPlaylists = originalPlaylists.filter(playlist => playlist.name.toLowerCase().includes(lowerValue) || (playlist.description && playlist.description.toLowerCase().includes(lowerValue)));
      }
      else{
        // restore full list
        filteredPlaylists = originalPlaylists;
      }
    });
  });

  onDestroy(() => {
    if (searchQueryUnsubscribe) searchQueryUnsubscribe(); // stop listening
    if (playlistsStoreUnsubscribe) playlistsStoreUnsubscribe(); // stop listening
  });
</script>

{#if filteredPlaylists.length > 0}
  <div class="row">
    {#each filteredPlaylists as playlist}
      <div class="col-12 col-lg-4 col-md-4 col-sm-6 mt-2">
        <PlaylistComponent {playlist} />
      </div>
    {/each}
  </div>
{:else}
  <p class="text-center">Working.....</p>
{/if}
