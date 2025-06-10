<script>
</script>

<main>
  <div class="container" style="position:fixed; inset:0; max-width:100vw; max-height:100vh;">
    <header class="top-bar">
      <input type="text" placeholder="Search..." class="search-input" />
      <button class="search-btn">Search</button>
      <button class="add-btn">+</button>
    </header>
    <section class="content">
      <!-- Main content goes here -->
      <div>
        <h2>Placeholder Content</h2>
        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque euismod, nisi eu consectetur consectetur, nisl nisi consectetur nisi, eu consectetur nisl nisi euismod nisi.</p>
        <p>Aliquam erat volutpat. Suspendisse potenti. Etiam ac mauris lectus. Integer non dui vitae urna cursus tincidunt.</p>
        <p>Morbi nec magna ac sem euismod commodo. Proin dictum, lorem a facilisis pretium, enim erat dictum urna, nec dictum enim enim nec enim.</p>
        <p>Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.</p>
        <p>Curabitur euismod, nisi eu consectetur consectetur, nisl nisi consectetur nisi, eu consectetur nisl nisi euismod nisi.</p>
        <p>Aliquam erat volutpat. Suspendisse potenti. Etiam ac mauris lectus. Integer non dui vitae urna cursus tincidunt.</p>
        <p>Morbi nec magna ac sem euismod commodo. Proin dictum, lorem a facilisis pretium, enim erat dictum urna, nec dictum enim enim nec enim.</p>
        <p>Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.</p>
      </div>
    </section>
    <nav class="middle-bar">
      <img class="left-btn" src="https://picsum.photos/100" alt="404" />
      <button class="middle-btn">Juice World</button>
      <button class="right-btn">▶</button>
    </nav>
    <nav class="bottom-bar">
      <button class="bottom-btn">Search</button>
      <button class="bottom-btn">Library</button>
      <button class="bottom-btn">Settings</button>
    </nav>
  </div>
</main>

<style>
.container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  margin: 0;
  padding: 0;
}

.top-bar {
  position: sticky;
  top: 0;
  background: black;
  padding: 0.75rem;
  display: flex;
  gap: 0.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.03);
  z-index: 10;
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
  border-bottom: 2px solid white;
}

.content {
  flex: 1 1 auto;
  overflow-y: auto;
  padding: 1rem;
  background: #000000;
}

.middle-bar {
  position: sticky;
  bottom: 3.5rem; /* Height of bottom-bar + gap */
  display: flex;
  justify-content: space-around;
  padding: 0.5rem 0;
  z-index: 9;
  border: 2px solid #444;
  width: 96%;
  margin-left: auto;
  margin-right: auto;
}

.middle-btn {
  flex: 1 1 0;
  margin: 0 0.25rem;
  padding: 0.75rem 0;
  background: #333;
  border: none;
  color: white;
  border-radius: 6px;
  font-size: 1rem;
}

.bottom-bar {
  position: sticky;
  bottom: 0;
  background: black;
  display: flex;
  justify-content: space-around;
  padding: 0.5rem 0;
  box-shadow: 0 -2px 4px rgba(0,0,0,0.03);
  z-index: 10;
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
  border-top: 2px solid white;
}

.bottom-btn {
  flex: 1 1 0;
  margin: 0 0.25rem;
  padding: 0.75rem 0;
  background: #2a2020;
  border: none;
  color: white;
  border-radius: 6px;
  font-size: 1rem;
}

.search-input {
  flex: 1 1 0;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

.search-btn {
  padding: 0.5rem 0.5rem;
  background: #ff3e00;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
}

.left-btn {
  width: 3rem;
  background: transparent;
  border: none;
  cursor: pointer;
}

.add-btn {
  padding: 0.5rem 1rem;
  background: #ff3e00;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
}

@media (max-width: 600px) {
  .content {
    padding: 0.5rem;
  }
  .top-bar, .bottom-bar, .middle-bar {
    padding-left: 0.5rem;
    padding-right: 0.5rem;
  }
  /* .middle-bar {
    width: 98%;
  } */
}
</style>
