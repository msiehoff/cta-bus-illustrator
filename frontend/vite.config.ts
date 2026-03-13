import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  // MapLibre GL 5 uses native class fields. Without this, esbuild compiles them
  // into __publicField() helpers that are missing inside MapLibre's web workers.
  optimizeDeps: {
    esbuildOptions: { target: 'esnext' },
  },
  build: {
    target: 'esnext',
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
