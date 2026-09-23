import { ref } from 'vue'


const isDark = ref(false)

function applyTheme() {
  const root = document.documentElement
  if (isDark.value) {
    root.classList.add('p-dark')
  } else {
    root.classList.remove('p-dark')
  }
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const saved = localStorage.getItem('theme')
  if (saved) {
    isDark.value = saved === 'dark'
  } else {
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
}

function toggleTheme() {
  isDark.value = !isDark.value
  applyTheme()
}

export function useTheme() {
  return { isDark, initTheme, toggleTheme }
}