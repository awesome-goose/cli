import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Builds into ../public (served by the Go backend); the dev server
// proxies /api requests to the Go backend on :8080.
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: '../public',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
