import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 5005,
    proxy: {
      '/api': {
        target: 'http://localhost:5050',
        changeOrigin: true,
      }
    }
  }
});
