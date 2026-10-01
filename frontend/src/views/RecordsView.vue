<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue"
import { useRoute } from "vue-router"
import api, { startOutput, stopOutput, fetchOpcData } from "../services/api"
import { useValveStore } from "../stores/valve"

import AutoComplete from "primevue/autocomplete"
import Button from "primevue/button"
import { useToast } from "primevue/usetoast"
import { Line } from "vue-chartjs"
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Filler,
} from "chart.js"

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Filler,
)

const STEPS = 40
const route = useRoute()
const valveStore = useValveStore()
const toast = useToast()

const valve = computed(() => valveStore.selectedValve)

const allValves = ref([])
const query = ref("")
const suggestions = ref([])
const pendingValve = ref(null)

const controlStatus = ref("stopped")
const isRecording = ref(false)
const isStarting = ref(false)
const isStopping = ref(false)
const imageViewerVisible = ref(false)
const recordedSession = ref([])

let liveDataInterval = null
const liveLabels = ref(Array(STEPS).fill(""))
const livePressure = ref(Array(STEPS).fill(0))
const liveFlow = ref(Array(STEPS).fill(0))
const liveCommand = ref(Array(STEPS).fill(0))
const liveFeedback = ref(Array(STEPS).fill(0))

const currentValues = ref({
  pressure: "0.0",
  flow: "0.0",
  command: "0.0",
  feedback: "0.0",
})

async function loadValves() {
  try {
    const response = await api.get("/valves")
    allValves.value = Array.isArray(response.data) ? response.data : []
  } catch (error) {
    console.error("Failed to load valve list:", error)
  }
}

async function fetchValveById(id) {
  if (!id) return
  try {
    const response = await api.get(`/valves/${id}`)
    valveStore.selectValve(response.data)
  } catch (error) {
    console.error("Failed to fetch valve:", error)
  }
}

function search(event) {
  const q = String(event.query || "").trim().toLowerCase()
  suggestions.value = allValves.value.filter((item) =>
    (item.part_number && item.part_number.toLowerCase().includes(q)) ||
    (item.manufacturer && item.manufacturer.toLowerCase().includes(q))
  )
}

function onSelect(event) {
  pendingValve.value = event.value
}

function confirmShow() {
  if (!pendingValve.value) return
  valveStore.selectValve(pendingValve.value)
  toast.add({
    severity: "info",
    summary: "Valve Configured",
    detail: `${pendingValve.value.part_number} loaded into test bench.`,
    life: 2500,
  })
}

async function startValve() {
  if (!valve.value || isStarting.value) return
  isStarting.value = true
  try {
    await startOutput()
    controlStatus.value = "active"
    toast.add({ severity: "success", summary: "HPU Active", detail: "Output signal engaged.", life: 3000 })
  } catch (error) {
    toast.add({ severity: "error", summary: "Start Failed", detail: "Failed to engage HPU.", life: 4000 })
  } finally {
    isStarting.value = false
  }
}

async function stopValve() {
  if (!valve.value || isStopping.value) return
  isStopping.value = true
  try {
    await stopOutput()
    controlStatus.value = "stopped"
    if (isRecording.value) toggleRecord()
    toast.add({ severity: "warn", summary: "HPU Stopped", detail: "Output signal disengaged.", life: 3000 })
  } catch (error) {
    toast.add({ severity: "error", summary: "Stop Failed", detail: "Failed to disengage HPU.", life: 4000 })
  } finally {
    isStopping.value = false
  }
}

async function toggleRecord() {
  if (!valve.value) return
  isRecording.value = !isRecording.value

  if (isRecording.value) {
    recordedSession.value = []
    toast.add({ severity: "info", summary: "Recording Started", detail: "Logging telemetry data.", life: 2500 })
  } else {
    if (recordedSession.value.length > 0) {
      try {
        await api.post('/records', {
          valve_id: valve.value.id,
          data_payload: JSON.stringify(recordedSession.value)
        })
        toast.add({ severity: "success", summary: "Recording Saved", detail: `Saved ${recordedSession.value.length} samples to database.`, life: 3000 })
      } catch (error) {
        toast.add({ severity: "error", summary: "Save Failed", detail: "Could not save to database.", life: 3000 })
      }
    }
  }
}

function exportRecordedCSV() {
  if (!recordedSession.value.length) return
  const headers = "Timestamp,Pressure_bar,Flow_Lmin,Command,Feedback\n"
  const rows = recordedSession.value
    .map((r) => `${r.time},${r.pressure},${r.flow},${r.command},${r.feedback}`)
    .join("\n")

  const blob = new Blob([headers + rows], { type: "text/csv;charset=utf-8;" })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.setAttribute("href", url)
  link.setAttribute("download", `TestReport_${valve.value?.part_number || "Valve"}_${Date.now()}.csv`)
  link.click()
}

async function pollLiveData() {
  const now = new Date()
  const timeStr = `${String(now.getHours()).padStart(2, "0")}:${String(now.getMinutes()).padStart(2, "0")}:${String(now.getSeconds()).padStart(2, "0")}`

  let p = 0, f = 0, c = 0, fb = 0
  try {
    const res = await fetchOpcData()
    const d = res.data
    p = d.pressure != null ? Number(d.pressure) : 0
    f = d.flow != null ? Number(d.flow) : 0
    c = d.command != null ? Number(d.command) : 0
    fb = d.feedback != null ? Number(d.feedback) : 0
  } catch (e) {
  }

  currentValues.value.pressure = p.toFixed(1)
  currentValues.value.flow = f.toFixed(1)
  currentValues.value.command = c.toFixed(1)
  currentValues.value.feedback = fb.toFixed(1)

  liveLabels.value.shift()
  liveLabels.value.push(timeStr)

  livePressure.value.shift()
  livePressure.value.push(p)

  liveFlow.value.shift()
  liveFlow.value.push(f)

  liveCommand.value.shift()
  liveCommand.value.push(c)

  liveFeedback.value.shift()
  liveFeedback.value.push(fb)

  if (isRecording.value) {
    recordedSession.value.push({ time: timeStr, pressure: p, flow: f, command: c, feedback: fb })
  }
}

const pressureChartData = computed(() => ({
  labels: [...liveLabels.value],
  datasets: [{ label: "Pressure (bar)", data: [...livePressure.value], borderColor: "#00CCFF", backgroundColor: "rgba(0, 204, 255, 0.12)", fill: true, tension: 0.35, pointRadius: 0, borderWidth: 2 }]
}))

const flowChartData = computed(() => ({
  labels: [...liveLabels.value],
  datasets: [{ label: "Flow (L/min)", data: [...liveFlow.value], borderColor: "#10b981", backgroundColor: "rgba(16, 185, 129, 0.08)", fill: true, tension: 0.35, pointRadius: 0, borderWidth: 2 }]
}))

const commandChartData = computed(() => ({
  labels: [...liveLabels.value],
  datasets: [{ label: "Command", data: [...liveCommand.value], borderColor: "#f59e0b", backgroundColor: "rgba(245, 158, 11, 0.12)", fill: true, tension: 0.35, pointRadius: 0, borderWidth: 2 }]
}))

const feedbackChartData = computed(() => ({
  labels: [...liveLabels.value],
  datasets: [{ label: "Feedback", data: [...liveFeedback.value], borderColor: "#8b5cf6", backgroundColor: "rgba(139, 92, 246, 0.12)", fill: true, tension: 0.35, pointRadius: 0, borderWidth: 2 }]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: { duration: 0 },
  interaction: { intersect: false, mode: "index" },
  plugins: {
    legend: { display: true, position: "bottom", labels: { boxWidth: 8, usePointStyle: true, padding: 6, font: { size: 9 } } },
    tooltip: { enabled: true }
  },
  scales: {
    x: { display: true, type: "category", grid: { display: false }, ticks: { color: "#64748b", font: { size: 9 }, maxTicksLimit: 6 } },
    y: { beginAtZero: true, border: { display: false }, grid: { color: "rgba(148, 163, 184, 0.12)" }, ticks: { color: "#64748b", font: { size: 9 } } }
  }
}

const imageUrl = computed(() => {
  if (!valve.value?.image_path) return null
  if (valve.value.image_path.startsWith("http")) return valve.value.image_path
  const base = import.meta.env.VITE_API_URL || ""
  return `${base}${valve.value.image_path}`
})

onMounted(async () => {
  await loadValves()
  const targetId = route.params.id || route.query.id
  if (targetId) await fetchValveById(targetId)
  liveDataInterval = setInterval(pollLiveData, 1000)
})

onBeforeUnmount(() => {
  if (liveDataInterval) clearInterval(liveDataInterval)
})
</script>

<template>
  <div class="records-page">
    <div class="search-toolbar">
      <div class="search-section">
        <AutoComplete
          v-model="query"
          :suggestions="suggestions"
          optionLabel="part_number"
          placeholder="Select Valve for Testing..."
          class="valve-search"
          @complete="search"
          @item-select="onSelect"
        />
        <Button label="Load" icon="pi pi-check" :disabled="!pendingValve" @click="confirmShow" />
      </div>

      <!-- Real-time Status & Telemetry Bar -->
      <div class="live-values-section">
        <div class="hpu-badge" :class="controlStatus">
          <span class="pulse-dot"></span>
          <span>HPU: {{ controlStatus.toUpperCase() }}</span>
        </div>
        <div class="live-value-item">
          <span class="label">Command</span>
          <span class="value">{{ currentValues.command }}</span>
        </div>
        <div class="live-value-item">
          <span class="label">Feedback</span>
          <span class="value">{{ currentValues.feedback }}</span>
        </div>
        <div class="live-value-item">
          <span class="label">Pressure</span>
          <span class="value">{{ currentValues.pressure }} <small>bar</small></span>
        </div>
        <div class="live-value-item">
          <span class="label">Flow</span>
          <span class="value">{{ currentValues.flow }} <small>L/min</small></span>
        </div>
      </div>
    </div>

    <div class="wireframe-grid">
      <!-- Left Controls Column -->
      <div class="sidebar-column">
        <section class="panel control-panel">
          <Button label="Start" icon="pi pi-play" severity="success" class="ctrl-btn" :loading="isStarting" :disabled="!valve || isStarting || controlStatus === 'active'" @click="startValve" />
          <Button label="Stop" icon="pi pi-power-off" severity="danger" class="ctrl-btn" :loading="isStopping" :disabled="!valve || isStopping || controlStatus === 'stopped'" @click="stopValve" />
          <Button :label="isRecording ? 'Halt Rec' : 'Record'" icon="pi pi-circle-fill" :severity="isRecording ? 'warn' : 'secondary'" class="ctrl-btn" :disabled="!valve" @click="toggleRecord" />
        </section>

        <!-- Valve Preview -->
        <section class="panel picture-panel">
          <header class="panel-header">
            <h3>Valve Preview</h3>
          </header>
          <div class="picture-content">
            <img v-if="imageUrl" :src="imageUrl" class="valve-image" alt="Valve Preview" />
            <div v-else class="image-placeholder">
              <i class="pi pi-image"></i>
              <span>{{ valve ? 'No Image Uploaded' : 'No Valve Loaded' }}</span>
            </div>
          </div>
        </section>

        <!-- Rating Spec Sheet -->
        <section class="panel info-panel">
          <div class="info-row">
            <span class="info-label">Part Number</span>
            <span class="info-val">{{ valve?.part_number || "-" }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Series / Type</span>
            <span class="info-val">{{ valve?.component_series || "-" }} / {{ valve?.valve_type || "-" }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Rated Flow</span>
            <span class="info-val">{{ valve?.rated_flow ?? "-" }} L/min</span>
          </div>
          <div class="info-row">
            <span class="info-label">Max Pressure</span>
            <span class="info-val">{{ valve?.max_pressure ?? "-" }} bar</span>
          </div>
          <div v-if="recordedSession.length > 0 && !isRecording" class="export-block">
            <Button label="Export CSV" icon="pi pi-download" severity="help" class="export-btn" @click="exportRecordedCSV" />
          </div>
        </section>
      </div>

      <!-- Graph Matrix -->
      <section class="panel graph-panel pressure-graph">
        <header class="panel-header"><h3>Pressure Response</h3></header>
        <div class="chart-container">
          <Line :data="pressureChartData" :options="chartOptions" />
        </div>
      </section>

      <section class="panel graph-panel command-graph">
        <header class="panel-header"><h3>Command Input</h3></header>
        <div class="chart-container">
          <Line :data="commandChartData" :options="chartOptions" />
        </div>
      </section>

      <section class="panel graph-panel flow-graph">
        <header class="panel-header"><h3>Flow Response</h3></header>
        <div class="chart-container">
          <Line :data="flowChartData" :options="chartOptions" />
        </div>
      </section>

      <section class="panel graph-panel feedback-graph">
        <header class="panel-header"><h3>Position Feedback</h3></header>
        <div class="chart-container">
          <Line :data="feedbackChartData" :options="chartOptions" />
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.records-page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  overflow: hidden;
}
.search-toolbar {
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.35rem 0.85rem;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  gap: 1rem;
}
.search-section {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  max-width: 440px;
  flex: 1;
}
.valve-search {
  flex: 1;
}
.valve-search :deep(.p-inputtext) {
  width: 100%;
}
.live-values-section {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}
.hpu-badge {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.25rem 0.6rem;
  border-radius: 4px;
  font-size: 0.72rem;
  font-weight: 700;
  font-family: Consolas, Monaco, monospace;
}
.hpu-badge.active {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}
.hpu-badge.stopped {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}
.pulse-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}
.live-value-item {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  min-width: 4.5rem;
}
.live-value-item .label {
  font-size: 0.65rem;
  color: var(--text-muted);
  text-transform: uppercase;
}
.live-value-item .value {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--primary-color);
  font-family: Consolas, Monaco, monospace;
  font-variant-numeric: tabular-nums;
}
.wireframe-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr) minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr) minmax(0, 1fr);
  grid-template-areas:
    "sidebar pressure flow"
    "sidebar command feedback";
  gap: 0.85rem;
}
.sidebar-column {
  grid-area: sidebar;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  min-height: 0;
}
.control-panel {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.4rem;
  padding: 0.6rem;
}
.ctrl-btn {
  width: 100%;
  padding: 0.45rem 0;
}
.picture-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 180px;
}
.picture-content {
  flex: 1;
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.valve-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.4rem;
  color: var(--text-muted);
  font-size: 0.8rem;
}
.image-placeholder i {
  font-size: 2rem;
}
.info-panel {
  padding: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}
.info-row {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid rgba(148, 163, 184, 0.15);
  padding-bottom: 0.3rem;
}
.info-label {
  font-size: 0.75rem;
  color: var(--text-muted);
}
.info-val {
  font-size: 0.85rem;
  font-weight: 600;
  font-family: Consolas, Monaco, monospace;
}
.export-block {
  margin-top: 0.4rem;
}
.export-btn {
  width: 100%;
}
.pressure-graph { grid-area: pressure; }
.flow-graph { grid-area: flow; }
.command-graph { grid-area: command; }
.feedback-graph { grid-area: feedback; }
.graph-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.chart-container {
  flex: 1;
  min-height: 0;
  padding: 0.4rem 0.6rem 0.1rem;
}
.panel {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}
.panel-header {
  padding: 0.4rem 0.75rem;
  background: rgba(148, 163, 184, 0.05);
  border-bottom: 1px solid var(--border-color);
}
.panel-header h3 {
  margin: 0;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-color);
}
</style>