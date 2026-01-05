import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
  optimizeDeps: {
    exclude: ["monaco-editor"], // 不预构建 Monaco Editor
  },
  base: "./",
  server: {
    fs: {
      strict: false,
    },
  },
  define: {
    "process.env.NODE_ENV": JSON.stringify(
      process.env.NODE_ENV || "development",
    ),
  },
  worker: {
    format: "es", // Worker 使用 ES 模块格式
  },
});
