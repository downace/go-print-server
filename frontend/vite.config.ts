import * as path from "node:path";
import { quasar, transformAssetUrls } from "@quasar/vite-plugin";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: { transformAssetUrls },
    }),

    // @quasar/plugin-vite options list:
    // https://github.com/quasarframework/quasar/blob/dev/vite-plugin/index.d.ts
    quasar(),
  ],
  resolve: {
    alias: {
      "@/runtime": path.resolve(__dirname, "./wailsjs/runtime"),
      "@/go": path.resolve(__dirname, "./wailsjs/go"),
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
