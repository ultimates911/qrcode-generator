<template>
  <div id="app">
    <router-view />
    <ThemeToggle />
  </div>
</template>

<script setup>
import ThemeToggle from '@/components/ThemeToggle.vue'
</script>

<style>

:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafb;
  --bg-tertiary: #f3f5f7;
  --text-primary: #2e3a44;
  --text-secondary: #666;
  --text-muted: #6c8c8c;
  --border-color: #edf0f2;
  --border-hover: #dfe6eb;
  --accent-color: #6c8c8c;
  --accent-hover: #5a7a7a;
  --error-color: #c54747;
  --shadow-light: rgba(58, 106, 149, 0.15);
  --shadow-dark: rgba(0, 0, 0, 0.1);
  --modal-overlay: rgba(44, 44, 44, 0.6);
}


:root.dark-theme {
  --bg-primary: #1a1a1a;
  --bg-secondary: #2d2d2d;
  --bg-tertiary: #3d3d3d;
  --text-primary: #ffffff;
  --text-secondary: #cccccc;
  --text-muted: #aaaaaa;
  --border-color: #404040;
  --border-hover: #505050;
  --accent-color: #7c9c9c;
  --accent-hover: #8cacac;
  --error-color: #e57373;
  --shadow-light: rgba(0, 0, 0, 0.3);
  --shadow-dark: rgba(0, 0, 0, 0.5);
  --modal-overlay: rgba(0, 0, 0, 0.7);
}

body {
  margin: 0;
  font-family: system-ui, sans-serif;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: background-color 0.3s ease, color 0.3s ease;
}

* {
  transition: background-color 0.3s ease, border-color 0.3s ease, color 0.3s ease;
}
</style>