<script setup>
import { ref } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import Topbar from '../components/Topbar.vue'

const sidebarCollapsed = ref(true)
</script>

<template>
  <div class="main-layout">
    <Sidebar :collapsed="sidebarCollapsed" />
    <div class="content-area">
      <Topbar @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed" />
      <main class="page-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.main-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}
.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 100vh;
  overflow: hidden;
}
.page-content {
  flex: 1;
  padding: 1.5rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
</style>