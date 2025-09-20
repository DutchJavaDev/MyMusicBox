<script>
  // @ts-nocheck
  import { searchQuery } from "../scripts/util";
</script>

<div class="search" role="search" aria-label="Search music">
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
    <path stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35" />
    <circle cx="11" cy="11" r="6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
  </svg>
  <input
    type="search"
    on:keydown={(e) => {
      searchQuery.set(e.target.value);
    }}
    placeholder="Search music #not working jet..."
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
    width: 95%;
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