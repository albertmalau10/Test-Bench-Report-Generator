<template>
  <div class="dashboard-view">
    <!-- Page header (fixed height) -->
    <div class="dashboard-header">
      <div>
        <h2 class="m-0"></h2>
        <p class="text-color-secondary m-0"></p>
      </div>
      <div class="header-actions">
        <Button label="Valve List" icon="pi pi-list" @click="goToValveList" />
        <Button label="Records" icon="pi pi-database" @click="goToRecords" />
      </div>
    </div>

    <!-- 1. Summary KPI / Header Cards (compact, 3 in a row) -->
    <div class="kpi-row">
      <Card class="kpi-card">
        <template #content>
          <div class="kpi-content">
            <div>
              <div class="kpi-label">Total Registered Valves</div>
              <div class="kpi-value">{{ summary.totalRegistered }}</div>
            </div>
            <i class="pi pi-sitemap kpi-icon text-primary"></i>
          </div>
        </template>
      </Card>

      <Card class="kpi-card">
        <template #content>
          <div class="kpi-content">
            <div>
              <div class="kpi-label">Active / Tested Valves</div>
              <div class="kpi-value">{{ summary.activeTested }}</div>
            </div>
            <i class="pi pi-check-circle kpi-icon text-green-500"></i>
          </div>
        </template>
      </Card>

      <Card class="kpi-card">
        <template #content>
          <div class="kpi-content">
            <div>
              <div class="kpi-label">Incomplete Data</div>
              <div class="kpi-value">{{ summary.incompleteData }}</div>
            </div>
            <i class="pi pi-exclamation-triangle kpi-icon text-orange-500"></i>
          </div>
        </template>
      </Card>
    </div>

    <!-- 2. Visualisasi Data & Spesifikasi -->
    <div class="dashboard-body">
      <Card class="dashboard-card">
        <template #title>Comparison of Valve Series</template>
        <template #content>
          <div class="chart-wrap">
            <Chart
              v-if="chartData.labels.length"
              type="pie"
              :data="chartData"
              :options="chartOptions"
              class="w-full h-full"
            />
            <div v-else class="empty-state">
              <i class="pi pi-chart-pie text-4xl text-color-secondary mb-2"></i>
              <div class="text-color-secondary">No data yet</div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

// PrimeVue components (explicit import in case they're not registered globally)
import Card from 'primevue/card'
import Button from 'primevue/button'
import Chart from 'primevue/chart'

// Base URL for your Go/Gin backend
const API_BASE = 'http://localhost:8080'

const router = useRouter()

const summary = ref({
  totalRegistered: 0,
  activeTested: 0, 
  incompleteData: 0,
})

const loadingSummary = ref(false)
const loadingChart = ref(false)

const chartData = ref({
  labels: [],
  datasets: [
    {
      data: [],
      backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#EF5350', '#AB47BC', '#26C6DA', '#8D6E63'],
    },
  ],
})

const chartOptions = ref({
  plugins: {
    legend: {
      position: 'right',
    },
  },
  maintainAspectRatio: false,
})

function goToValveList() {
  router.push('/valves')
}

function goToRecords() {
  router.push('/records')
}

function isIncomplete(valve) {
  return (
    !valve.image_path ||
    !valve.datasheet_path ||
    !valve.manufacturer ||
    !valve.part_number ||
    !valve.component_series ||
    !valve.valve_type ||
    !valve.command_type ||
    valve.command_value === null ||
    valve.command_value === undefined ||
    valve.command_value === ''
  )
}

async function fetchSummary(valves) {
  loadingSummary.value = true
  try {
    summary.value.totalRegistered = valves.length
    summary.value.incompleteData = valves.filter(isIncomplete).length
    // activeTested left at 0 for now — depends on the Records endpoint,
    // which we're wiring up separately.
  } finally {
    loadingSummary.value = false
  }
}

async function fetchChartData(valves) {
  loadingChart.value = true
  try {
    // Group valves by component_series (e.g. "3X", "4X", "5X")
    const counts = {}
    for (const valve of valves) {
      const key = valve.component_series || 'Unspecified'
      counts[key] = (counts[key] || 0) + 1
    }
    chartData.value.labels = Object.keys(counts)
    chartData.value.datasets[0].data = Object.values(counts)
  } finally {
    loadingChart.value = false
  }
}

async function loadDashboardData() {
  try {
    const { data } = await axios.get(`${API_BASE}/valves`)
    const valves = Array.isArray(data) ? data : []
    await fetchSummary(valves)
    await fetchChartData(valves)
  } catch (err) {
    console.error('Failed to load dashboard data:', err)
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
/* Fills the height given by the layout. Uses overflow-y: auto as a safety
   net — if content ever ends up taller than the available space (e.g. on a
   very short window), it scrolls internally instead of getting clipped. */
.dashboard-view {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 1rem;
  overflow-y: auto;
  overflow-x: hidden;
  box-sizing: border-box;
}

.dashboard-header {
  flex: 0 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

/* KPI row: plain CSS grid, doesn't depend on PrimeFlex being installed */
.kpi-row {
  flex: 0 0 auto;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

@media (max-width: 768px) {
  .kpi-row {
    grid-template-columns: 1fr;
  }
}

/* Compact KPI cards, similar sizing to ValveDetailsView cards */
.kpi-card :deep(.p-card-body) {
  padding: 0.85rem 1rem;
}

.kpi-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.kpi-label {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  margin-bottom: 0.25rem;
}

.kpi-value {
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1.2;
}

.kpi-icon {
  font-size: 1.5rem;
}

/* Chart section: takes the remaining space */
.dashboard-body {
  flex: 1 1 auto;
  min-height: 260px; /* guarantees the chart card has room even in a short viewport */
  display: flex;
}

.dashboard-card {
  display: flex;
  flex-direction: column;
  width: 100%;
}

/* PrimeVue Card content area should stretch and allow its child to shrink */
.dashboard-card :deep(.p-card-body),
.dashboard-card :deep(.p-card-content) {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.chart-wrap {
  flex: 1 1 auto;
  min-height: 220px;
  position: relative;
}

.empty-state {
  flex: 1 1 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>