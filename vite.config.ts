import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      assets: '/app/assets',
      pages: '/app/pages',
      shared: '/app/shared',
      app: '/app',
    },
  },
})
