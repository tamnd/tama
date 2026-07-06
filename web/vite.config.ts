import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// In dev, Vite serves the app and proxies the API to a running `tama serve`.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:4321',
    },
  },
})
