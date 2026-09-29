<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { API_ORIGIN } from '../services/api'
import { useRoute } from 'vue-router'
import api from '../services/api'
import { useValveStore } from '../stores/valve'
import AutoComplete from 'primevue/autocomplete'
import Button from 'primevue/button'

const route = useRoute()
const valveStore = useValveStore()
const valve = computed(() => valveStore.selectedValve)

const allValves = ref([])
const query = ref('')
const suggestions = ref([])
const pendingValve = ref(null)

async function loadValves() {
  const res = await api.get('/valves')
  allValves.value = res.data
}

async function fetchValveById(id) {
  if (!id) return
  try {
    const res = await api.get(`/valves/${id}`)
    valveStore.selectValve(res.data)
  } catch (error) {
    console.error('Failed to get Valve details:', error)
  }
}

function search(event) {
  const q = event.query.toLowerCase()
  suggestions.value = allValves.value.filter(v =>
    v.part_number?.toLowerCase().includes(q)
  )
}

function onSelect(event) {
  pendingValve.value = event.value
}

function confirmShow() {
  if (pendingValve.value) {
    valveStore.selectValve(pendingValve.value)
  }
}

function openDatasheet() {
  if (!datasheetUrl.value) return
  window.open(datasheetUrl.value, '_blank', 'noopener')
}

const imageUrl = computed(() => {
  if (!valve.value?.image_path) return null
  if (valve.value.image_path.startsWith('http')) return valve.value.image_path
<<<<<<< HEAD
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${baseUrl}${valve.value.image_path}`
=======
  return `${API_ORIGIN}${valve.value.image_path}`
>>>>>>> pr-theme-update
})

const datasheetUrl = computed(() => {
  if (!valve.value?.datasheet_path) return null
  if (valve.value.datasheet_path.startsWith('http')) return valve.value.datasheet_path
<<<<<<< HEAD
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${baseUrl}${valve.value.datasheet_path}`
=======
  return `${API_ORIGIN}${valve.value.datasheet_path}`
>>>>>>> pr-theme-update
})

const feedback = {
  position: 65,
  current: 1.25,
  temperature: 42,
  status: 'Normal'
}

onMounted(() => {
  loadValves()
  if (route.params.id) {
    fetchValveById(route.params.id)
  }
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    fetchValveById(newId)
  }
})
</script>

<template>
  <div class="records-page">
    <div class="search-bar">
      <AutoComplete
        v-model="query"
        :suggestions="suggestions"
        optionLabel="part_number"
        placeholder="Search Valve"
        @complete="search"
        @item-select="onSelect"
        style="flex: 1"
      />
      <Button label="Show Valve" icon="pi pi-check" @click="confirmShow" :disabled="!pendingValve" />
    </div>

    <div v-if="!valve" class="empty-state">
      <i class="pi pi-info-circle"></i>
      <p>Find and select the Valve above, then click Show Valve</p>
    </div>

    <div v-else class="grid">
      <div class="card">
        <h3>Picture of Valve</h3>
        <img v-if="imageUrl" :src="imageUrl" alt="valve" class="valve-img" />
        <div v-else class="valve-img placeholder"><i class="pi pi-image"></i></div>
        <p class="highlight">{{ valve.part_number }}</p>
        <p class="muted">{{ valve.component_series }} {{ valve.valve_type }}</p>
      </div>

      <div class="card">
        <h3>Pressure</h3>
        <p class="big-number"><strong>Max Pressure :</strong> {{ valve.max_pressure ?? '-' }} Bar</p>
      </div>

      <div class="card">
        <h3>Flow</h3>
        <p><strong>Rated Flow:</strong> {{ valve.rated_flow ?? '-' }}</p>
        <p><strong>Max Flow:</strong> {{ valve.max_flow ?? '-' }}</p>
      </div>

      <div class="card">
        <h3>Valve Datasheet</h3>
        <div v-if="datasheetUrl" class="datasheet-box">
          <i class="pi pi-file-pdf datasheet-icon"></i>
          <Button label="View Datasheet" icon="pi pi-external-link" @click="openDatasheet" />
        </div>
        <p v-else class="muted">There is no datasheet available for this valve</p>
      </div>

      <div class="card">
        <h3>Command</h3>
        <p><strong>Type:</strong> {{ valve.command_type || '-' }}</p>
        <p><strong>Value:</strong> {{ valve.command_value ?? '-' }}</p>
      </div>

      <div class="card">
        <h3>Feedback</h3>
        <p><strong>Position:</strong> {{ feedback.position }}%</p>
        <p><strong>Current:</strong> {{ feedback.current }} A</p>
        <p><strong>Temperature:</strong> {{ feedback.temperature }} °C</p>
        <p><strong>Status:</strong> {{ feedback.status }}</p>
        <small class="muted">*simulation data, not connected in real-time</small>
      </div>
    </div>
  </div>
</template>

<style scoped>
.records-page {
  height: calc(100vh - 30px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.search-bar {
  display: flex;
  gap: 0.75rem;
  margin: 0 0 1rem 0;
  flex-shrink: 0;
}
.grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(2, 1fr);
  gap: 1rem;
}
.card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 1.25rem 1.5rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.datasheet-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 0;
}
.datasheet-icon {
  font-size: 2.5rem;
  color: #dc2626;
}
.card h3 { margin: 0 0 0.75rem 0; font-size: 1.05rem; color: var(--text-muted); }
.big-number { font-size: 2.25rem; font-weight: 400; color: #ffffff; margin: 0; }
.unit { font-size: 1rem; color: var(--text-muted); }
.highlight { font-weight: 600; color: var(--primary-color); margin: 0.5rem 0 0 0; font-size: 1.1rem; }
.muted { color: var(--text-muted); font-size: 0.9rem; margin: 0.25rem 0 0 0; }
.card p { margin: 0.35rem 0; font-size: 1rem; }
.valve-img { width: 100%; height: 130px; object-fit: contain; background: var(--bg-color); border-radius: 8px; }
.valve-img.placeholder { display: flex; align-items: center; justify-content: center; font-size: 2rem; color: var(--text-muted); }
.empty-state { text-align: center; padding: 2rem 1rem; color: var(--text-muted); }
.empty-state i { font-size: 1.75rem; display: block; margin-bottom: 0.5rem; }
</style>