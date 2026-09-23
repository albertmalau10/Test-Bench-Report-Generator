<!-- frontend/src/views/SettingView.vue -->
<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const loading = ref(false)
const generating = ref(false)

const settings = ref({
  modbus_server_address: '',
  ctrlx_host: '',
  ctrlx_username: '',
  ctrlx_password: ''
})

async function loadSettings() {
  try {
    const { data } = await api.get('/settings')
    settings.value.modbus_server_address = data.modbus_server_address || ''
    settings.value.ctrlx_host = data.ctrlx_host || ''
    settings.value.ctrlx_username = data.ctrlx_username || ''
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load settings', life: 3000 })
  }
}

async function saveSettings() {
  loading.value = true
  try {
    await api.put('/settings', settings.value)
    toast.add({ severity: 'success', summary: 'Success', detail: 'Settings saved', life: 3000 })
    settings.value.ctrlx_password = '' // Clear password field after save
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save settings', life: 3000 })
  } finally {
    loading.value = false
  }
}

async function generateToken() {
  generating.value = true
  try {
    await api.post('/settings/ctrlx/token')
    toast.add({ severity: 'success', summary: 'Token Generated', detail: 'ctrlX API token updated successfully', life: 3000 })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.error || 'Failed to generate token', life: 4000 })
  } finally {
    generating.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="settings-page">
    <Toast />
    <h2>System Settings</h2>
    
    <div class="settings-card">
      <div class="field">
        <label>Modbus Server Address</label>
        <InputText v-model="settings.modbus_server_address" placeholder="e.g., 127.0.0.1:502" />
      </div>

      <div class="field">
        <label>ctrlX Host</label>
        <InputText v-model="settings.ctrlx_host" placeholder="e.g., 192.168.1.1:8443" />
      </div>

      <div class="field">
        <label>ctrlX Username</label>
        <InputText v-model="settings.ctrlx_username" />
      </div>

      <div class="field">
        <label>ctrlX Password (Update Only)</label>
        <Password v-model="settings.ctrlx_password" :feedback="false" toggleMask />
      </div>

      <div class="actions">
        <Button label="Save Settings" icon="pi pi-save" :loading="loading" @click="saveSettings" />
        <Button label="Generate ctrlX Token" icon="pi pi-key" severity="secondary" :loading="generating" @click="generateToken" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  padding: 1rem;
}
.settings-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 2rem;
  max-width: 600px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}
</style>