<script>
  import { onMount } from "svelte";
  import PlaylistComponent from "../components/PlaylistComponent.svelte";
  import { playlistsStore } from "../scripts/playlistService.js";

  $: $playlistsStore;

  onMount(() => {});
</script>
{#if $playlistsStore.length > 0}
<div class="row">
   {#each $playlistsStore as playlist}
      <div class="col-12 col-lg-4 col-md-4 col-sm-6 mt-2">
        <PlaylistComponent {playlist} />
      </div>
    {/each}
</div>
{:else}
  <p class="text-center">Working.....</p>
{/if}