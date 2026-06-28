import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      "/shorten": {
        target: "http://backend:8080",
      },
      "/stats": {
        target: "http://backend:8080",
      },
    },
  },
});
