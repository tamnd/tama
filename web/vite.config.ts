import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// In dev, Vite serves the app and proxies the API to a running `tama serve`.
export default defineConfig({
  plugins: [react()],
  define: {
    // Ships the /dev/gallery route in a production build when GALLERY=1.
    __GALLERY__: JSON.stringify(process.env.GALLERY === '1'),
  },
  server: {
    proxy: {
      '/api': 'http://localhost:4321',
    },
  },
})
