<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

defineEmits(['toggle-sidebar'])

const authStore = useAuthStore()

const roleLabel = computed(() => {
  const role = authStore.user?.role
  if (!role) return ''
  return role.charAt(0).toUpperCase() + role.slice(1)
})
</script>

<template>
  <header class="topbar">
    <button class="icon-btn" @click="$emit('toggle-sidebar')">
      <i class="pi pi-bars"></i>
    </button>

    <div class="topbar-right">
      <button class="icon-btn">
        <i class="pi pi-bell"></i>
      </button>
      <div class="profile">
        <i class="pi pi-user-circle"></i>
        <span>{{ roleLabel }}</span>
      </div>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  height: 64px;
  min-height: 64px;
  max-height: 64px;
  flex-shrink: 0;
  flex-grow: 0;
  box-sizing: border-box;
  background: var(--topbar-bg);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.25rem;
  position: sticky;
  top: 0;
  z-index: 10;
}
.icon-btn {
  background: none;
  border: none;
  font-size: 1.1rem;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0.4rem;
  border-radius: 6px;
}
.icon-btn:hover {
  background: var(--sidebar-active-bg);
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.profile {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: var(--text-color);
  font-weight: 500;
}
.profile i {
  font-size: 1.6rem;
}
</style>