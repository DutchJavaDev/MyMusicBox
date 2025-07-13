<script>
  import { onMount } from "svelte";
  import PlaylistComponent from "../components/PlaylistComponent.svelte";
  import { playlistsStore } from "../scripts/playlistService.js";

  $: $playlistsStore;

  onMount(() => {});
</script>
{#if $playlistsStore.length > 0}
  {#each $playlistsStore as playlist}
    <PlaylistComponent {playlist} />
  {/each}
{:else}
  <p class="text-center">Working.....</p>
{/if}