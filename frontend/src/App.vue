<script setup>
import { onMounted } from 'vue'
import Toast from 'primevue/toast'
import { useTheme } from './composables/useTheme'

const { initTheme } = useTheme()
onMounted(() => {
  initTheme()
})
</script>

<template>
  <div class="app-container">
    <Toast />
    <router-view />
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  overflow: hidden;
}
</style>