import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// VITE_BASE_PATH is set by packaging/build-frontend.sh to
// "/valve-database-app/" for the ctrlX CORE build, so every asset URL
// and Vue Router link resolves correctly behind ctrlX's reverse
// proxy. Local `npm run dev` / `npm run build` without that env var
// keeps the default "/".
// https://vite.dev/config/
export default defineConfig({
  base: process.env.VITE_BASE_PATH || '/',
  plugins: [vue()],
  server: {
    // api.js calls relative "<BASE_URL>api/..." paths (e.g. "/api/valves")
    // so the same code works unchanged behind ctrlX CORE's reverse proxy.
    // This proxy makes that work for `npm run dev` too, forwarding to the
    // Go backend (`go run .` in backend/, default :8080) without needing
    // VITE_API_URL or backend CORS at all.
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
