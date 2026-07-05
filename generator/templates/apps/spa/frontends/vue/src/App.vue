<script setup>
import { onMounted, ref } from 'vue'

const info = ref(null)
const error = ref(null)

onMounted(async () => {
  try {
    const res = await fetch('/api/version')
    info.value = await res.json()
  } catch (err) {
    error.value = String(err)
  }
})
</script>

<template>
  <main class="app">
    <h1>Goose + Vue</h1>
    <p>
      This page is served by the Goose SPA platform. The backend answers
      JSON under <code>/api</code>.
    </p>
    <p v-if="info" class="status">
      Backend says: <strong>{{ info.name }}</strong> v{{ info.version }}
    </p>
    <p v-if="error" class="error">Backend unreachable: {{ error }}</p>
  </main>
</template>
