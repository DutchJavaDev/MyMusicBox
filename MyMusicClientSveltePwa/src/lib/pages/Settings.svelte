<script>
  import { onMount } from "svelte";
  import { clearStorage, getConfiguration, setConfiguration } from "../scripts/storageService";

  function handleSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const sleepTimer = formData.get("sleepTimer");
    const fetchTimer = formData.get("fetchTimer");
    const byPassCache = formData.get("byPassCache") === "on";
    
    configuarion.sleepTimer = sleepTimer;
    configuarion.fetchTimer = fetchTimer;
    configuarion.byPassCache = byPassCache;

    setConfiguration(configuarion);
  }

  let configuarion;
  let sleepTimer = $state(0);
  let fetchTimer = $state(0);
  let byPassCache = $state(false);

  onMount(() => {
    configuarion = getConfiguration();
    sleepTimer = configuarion.sleepTimer;
    fetchTimer = configuarion.fetchTimer;
    byPassCache = configuarion.byPassCache;
  });
</script>

<!-- svelte-ignore event_directive_deprecated -->
<form on:submit|preventDefault={handleSubmit} class="p-2 rounded rounded-2 tile-bg">
  <div class="mb-3">
    <label for="sleepTimer" class="form-label">Sleep Timer (Minutes)</label>
    <input type="number" required class="form-control" id="sleepTimer" name="sleepTimer" value={sleepTimer} placeholder="Stop music after some time" />
  </div>
  <div class="mb-3">
    <label for="fetchTimer" class="form-label">Fetch Timer (Seconds)</label>
    <input type="number" required class="form-control" id="fetchTimer" name="fetchTimer" value={fetchTimer} placeholder="Api interval time" />
  </div>
  <div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
    <label for="byPassCache" class="form-label">Ignore Cached urls</label>
    <input type="checkbox" checked={byPassCache} id="byPassCache" name="byPassCache" class="form-check-input" />
  </div>
  <button type="submit" class="btn btn-primary">Save Settings</button>
</form>

<div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
  <label for="localStorageClear" class="form-label">Local Storage !Expirmental</label>
  <!-- This is causing some issues, playback breaks and updating ui freezes -->
  <!-- svelte-ignore event_directive_deprecated -->
  <button id="localStorageClear" class="btn btn-danger" on:click={() => clearStorage()}>Clear Local Storage</button>
</div>

<div class="mt-3 mb-3 p-2 rounded rounded-2 tile-bg">
  <button id="exportData" class="btn btn-secondary">View Server Logs</button>
</div>

<style>
  label {
    font-weight: bold;
    font-size: 1.1rem;
  }

  .tile-bg {
    background-color: #2c2c2c;
  }
</style>
