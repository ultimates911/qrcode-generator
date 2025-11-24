<template>
  <div class="theme-switch" :class="{ 'dark': isDark }" @click="toggleTheme" :title="isDark ? 'Светлая тема' : 'Тёмная тема'">
    <div class="switch-track">
      <div class="switch-thumb">
      </div>
      <div class="icons">
        <span class="sun">☀️</span>
        <span class="moon">🌙</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useTheme } from '../composables/useTheme'

const { isDark, toggleTheme } = useTheme()
</script>

<style scoped>
.theme-switch {
  position: fixed;
  bottom: 20px;
  right: 20px;
  cursor: pointer;
  z-index: 1000;
}

.switch-track {
  position: relative;
  width: 75px;
  height: 32px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 
    0 4px 15px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  padding: 0 5px;
}

.switch-thumb {
  position: absolute;
  width: 22px; 
  height: 22px; 
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(5px);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  left: 12px; 
}

.theme-switch.dark .switch-thumb {
  left: 54px; 
  background: rgba(108, 140, 140, 0.9);
}

.icons {
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
}

.sun, .moon {
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 1; 
}

.theme-switch:hover .switch-track {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
  box-shadow: 
    0 6px 20px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.theme-switch:hover .switch-thumb {
  transform: scale(1.05);
  box-shadow: 
    0 3px 12px rgba(0, 0, 0, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.theme-switch:active .switch-thumb {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .theme-switch {
    bottom: 15px;
    right: 15px;
  }
  
  .switch-track {
    width: 65px;
    height: 30px;
  }
  
  .switch-thumb {
    width: 20px; 
    height: 20px; 
    left: 10px; 
  }
  
  .theme-switch.dark .switch-thumb {
    left: 35px; 
  }
}
</style>