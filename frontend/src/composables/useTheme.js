import { ref, onMounted, watch } from 'vue'

export function useTheme() {
  const isDark = ref(false)

  onMounted(() => {
    const saved = localStorage.getItem('theme')
    const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    
    isDark.value = saved ? saved === 'dark' : systemDark
    applyTheme()
  })

  function toggleTheme() {
    isDark.value = !isDark.value
    applyTheme()
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  }

  function applyTheme() {
    if (isDark.value) {
      document.documentElement.classList.add('dark-theme')
    } else {
      document.documentElement.classList.remove('dark-theme')
    }
  }

  return {
    isDark,
    toggleTheme
  }
}