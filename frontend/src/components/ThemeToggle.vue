<template>
  <button class="theme-toggle" @click="toggleTheme" :title="isDark ? 'Светлая тема' : 'Тёмная тема'">
    <span class="theme-icon">{{ isDark ? '☀️' : '🌙' }}</span>
  </button>
</template>

<script setup>
import { useTheme } from '@/composables/useTheme'

const { isDark, toggleTheme } = useTheme()
</script>

<style scoped>
.theme-toggle {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  z-index: 1000;
  border: 1px solid var(--border-color);
}

.theme-toggle:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.theme-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>