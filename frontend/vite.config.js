import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/shorten': 'http://localhost:8080',
      '/stats':   'http://localhost:8080',
      '/health':  'http://localhost:8080',
    },
  },
  build: {
    outDir: '../static',
    emptyOutDir: true,
  },
})
