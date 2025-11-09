<script>
  // @ts-nocheck
  import { searchQuery } from "../scripts/util";
  import { getSearchQueryFromStorage, setSearchQueryInStorage } from "../scripts/storageService";
  import {  onMount } from "svelte";
  import { component } from "../scripts/routeService";
  
  let query = '';
  
  onMount(() => {
    // Initialize the search query from storage
    const storedQuery = getSearchQueryFromStorage();
    if (storedQuery && storedQuery.length > 0) {
      searchQuery.set(storedQuery);
      query = storedQuery;
    }

    // There is a x on right side of the input that triggers a 'search' event when clicked
    // it cleares the input but we also need to clear our stores and storage... :))))) javaScript
    document.getElementById('search-input').addEventListener('search', (e) => {
      searchQuery.set('');
      setSearchQueryInStorage('');
      query = '';
    });

      // update component on search query change
      component.subscribe(() => {
        searchQuery.set(getSearchQueryFromStorage());
      });
  });
</script>

<div class="search border border-1 border-dark" role="search" aria-label="Search music">
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
    <path stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35" />
    <circle cx="11" cy="11" r="6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
  </svg>
  <input 
    id="search-input"
    bind:value={query}
    type="search"
    on:keyup={(e) => {
      setSearchQueryInStorage(query);
      searchQuery.set(query);
    }}
    on:change={(e) => {
      console.log("change event", e);
      setSearchQueryInStorage(query);
      searchQuery.set(query);
    }}
    placeholder="Search and you shall find..."
    aria-label="Search"
  />
</div>
<style>
      .search {
    background: #2a2a2a;
    padding: 14px 18px;
    border-radius: 28px;
    display: flex;
    align-items: center;
    gap: 12px;
    box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.4);
    width: 95vw;
    margin: 0px auto;
  }
  .search svg {
    color: #b3b3b3;
    flex-shrink: 0;
  }
  .search input {
    background: transparent;
    border: 0;
    outline: 0;
    color: #ffffff;
    font-size: 16px;
    width: 100%;
  }
  .search input::placeholder {
    color: #b3b3b3;
  }
</style>