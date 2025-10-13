// vite.config.js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Можно задавать через переменные окружения Docker Compose:
//   VITE_BACKEND_HOST=http://qrcode-backend:8080
//   VITE_USE_POLLING=true
const backendTarget =
  process.env.VITE_BACKEND_HOST || 'http://qrcode-backend:8080'

const usePolling =
  String(process.env.CHOKIDAR_USEPOLLING || process.env.VITE_USE_POLLING || '')
    .toLowerCase() === 'true'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: true,          // слушать 0.0.0.0 внутри контейнера
    port: 5173,
    strictPort: true,
    open: false,
    watch: {
      usePolling,        // полезно на Mac/Windows/WSL
    },
    hmr: {
      // Когда порт опубликован как 5173:5173, этого достаточно.
      // Добавь clientPort: 5173, если браузер не цепляется к HMR за NAT/прокси.
      // clientPort: 5173,
    },
    proxy: {
      // Прокидываем API на бэкенд-сервис Docker по его имени
      // Меняй пути под своё API
      '/api': {
        target: backendTarget,
        changeOrigin: true,
        secure: false,
        // например, если бэкенд слушает на /, а у тебя /api на фронте:
        // rewrite: (p) => p.replace(/^\/api/, ''),
      },
    },
  },
})
