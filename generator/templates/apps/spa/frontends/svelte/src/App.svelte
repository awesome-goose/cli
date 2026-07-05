<script>
  let info = $state(null)
  let error = $state(null)

  $effect(() => {
    fetch('/api/version')
      .then((res) => res.json())
      .then((data) => (info = data))
      .catch((err) => (error = String(err)))
  })
</script>

<main class="app">
  <h1>Goose + Svelte</h1>
  <p>
    This page is served by the Goose SPA platform. The backend answers
    JSON under <code>/api</code>.
  </p>
  {#if info}
    <p class="status">
      Backend says: <strong>{info.name}</strong> v{info.version}
    </p>
  {/if}
  {#if error}
    <p class="error">Backend unreachable: {error}</p>
  {/if}
</main>
