<script></script>

{#each Array.from({ length: 100 }) as _, index}
  <div class="content-item">
    <!-- <Song {song} /> -->
    Content {index + 1}
  </div>
{/each}
<style>
    
</style>