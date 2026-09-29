<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import api, { startOutput, stopOutput, fetchOpcData } from "../services/api";
import { useValveStore } from "../stores/valve";

import AutoComplete from "primevue/autocomplete";
import Button from "primevue/button";
import { useToast } from "primevue/usetoast";

import { Line } from "vue-chartjs";
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
} from "chart.js";

// Register CategoryScale for the Time labels on the X-axis
ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Filler,
);

const STEPS = 40;

const route = useRoute();
const valveStore = useValveStore();
const toast = useToast();

const valve = computed(() => valveStore.selectedValve);

const allValves = ref([]);
const query = ref("");
const suggestions = ref([]);
const pendingValve = ref(null);

const controlStatus = ref("stopped");
const isRecording = ref(false);
const isStarting = ref(false);
const isStopping = ref(false);
const imageViewerVisible = ref(false);

// LIVE DATA STATE ARRAYS FOR GRAPHS
let liveDataInterval = null;
const liveLabels = ref(Array(STEPS).fill("")); // Array for Time Strings
const livePressureHistory = ref(Array(STEPS).fill(0));
const liveFlowHistory = ref(Array(STEPS).fill(0));
const liveCommandHistory = ref(Array(STEPS).fill(0));
const liveFeedbackHistory = ref(Array(STEPS).fill(0));

const currentValues = ref({
  pressure: "0.0",
  flow: "0.0",
  command: "0.0",
  feedback: "0.0",
});

async function loadValves() {
  try {
    const response = await api.get("/valves");
    allValves.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error("Failed to fetch valve list:", error);
    toast.add({
      severity: "error",
      summary: "Load Failed",
      detail: "Unable to load the valve list.",
      life: 3000,
    });
  }
}

async function fetchValveById(id) {
  if (!id) return;
  try {
    const response = await api.get(`/valves/${id}`);
    valveStore.selectValve(response.data);
  } catch (error) {
    console.error("Failed to fetch valve details:", error);
  }
}

function search(event) {
  const keyword = String(event.query || "")
    .trim()
    .toLowerCase();
  if (!keyword) {
    suggestions.value = allValves.value;
    return;
  }
  suggestions.value = allValves.value.filter((item) => {
    const partNumber = item.part_number?.toLowerCase() || "";
    const manufacturer = item.manufacturer?.toLowerCase() || "";
    return partNumber.includes(keyword) || manufacturer.includes(keyword);
  });
}

function onSelect(event) {
  pendingValve.value = event.value;
}

function confirmShow() {
  if (!pendingValve.value) return;
  valveStore.selectValve(pendingValve.value);
  toast.add({
    severity: "info",
    summary: "Valve Selected",
    detail: `${pendingValve.value.part_number} is ready for testing.`,
    life: 2500,
  });
}

async function startValve() {
  if (!valve.value || isStarting.value) return;
  isStarting.value = true;
  try {
    await startOutput();
    controlStatus.value = "active";
    toast.add({
      severity: "success",
      summary: "START",
      detail: "Output ON",
      life: 3000,
    });
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "START FAILED",
      detail: "Unable to activate the output.",
      life: 4000,
    });
  } finally {
    isStarting.value = false;
  }
}

async function stopValve() {
  if (!valve.value || isStopping.value) return;
  isStopping.value = true;
  try {
    await stopOutput();
    controlStatus.value = "stopped";
    isRecording.value = false;
    toast.add({
      severity: "warn",
      summary: "STOP",
      detail: "Output OFF",
      life: 3000,
    });
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "STOP FAILED",
      detail: "Unable to stop the output.",
      life: 4000,
    });
  } finally {
    isStopping.value = false;
  }
}

function toggleRecord() {
  if (!valve.value) return;
  isRecording.value = !isRecording.value;
  toast.add({
    severity: isRecording.value ? "info" : "warn",
    summary: isRecording.value ? "RECORDING STARTED" : "RECORDING STOPPED",
    detail: isRecording.value
      ? "Data recording active."
      : "Data recording stopped.",
    life: 3000,
  });
}

const imageUrl = computed(() => {
  if (!valve.value?.image_path) return null;
  if (valve.value.image_path.startsWith("http")) return valve.value.image_path;
  const baseUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
  return `${baseUrl}${valve.value.image_path}`;
});

function openImageViewer() {
  if (!imageUrl.value) return;
  imageViewerVisible.value = true;
  document.body.style.overflow = "hidden";
}

function closeImageViewer() {
  imageViewerVisible.value = false;
  document.body.style.overflow = "";
}

function handleImageViewerKeydown(event) {
  if (event.key === "Escape" && imageViewerVisible.value) {
    closeImageViewer();
  }
}

// FETCH LIVE OPC UA DATA
async function pollLiveData() {
  // Generate current timestamp string (HH:MM:SS)
  const now = new Date();
  const timeStr = `${now.getHours().toString().padStart(2, "0")}:${now.getMinutes().toString().padStart(2, "0")}:${now.getSeconds().toString().padStart(2, "0")}`;

  try {
    const response = await fetchOpcData();
    const data = response.data;

    currentValues.value.pressure =
      data.pressure !== null ? Number(data.pressure).toFixed(1) : "0.0";
    currentValues.value.flow =
      data.flow !== null ? Number(data.flow).toFixed(1) : "0.0";
    currentValues.value.command =
      data.command !== null ? Number(data.command).toFixed(1) : "0.0";
    currentValues.value.feedback =
      data.feedback !== null ? Number(data.feedback).toFixed(1) : "0.0";

    // Reassign completely new arrays to force Vue reactivity and Chart.js re-renders
    liveLabels.value = [...liveLabels.value.slice(1), timeStr];
    livePressureHistory.value = [
      ...livePressureHistory.value.slice(1),
      data.pressure !== null ? Number(data.pressure) : 0,
    ];
    liveFlowHistory.value = [
      ...liveFlowHistory.value.slice(1),
      data.flow !== null ? Number(data.flow) : 0,
    ];
    liveCommandHistory.value = [
      ...liveCommandHistory.value.slice(1),
      data.command !== null ? Number(data.command) : 0,
    ];
    liveFeedbackHistory.value = [
      ...liveFeedbackHistory.value.slice(1),
      data.feedback !== null ? Number(data.feedback) : 0,
    ];
  } catch (error) {
    // If the server disconnects, push 0s but keep the time moving forward to show a flatline
    liveLabels.value = [...liveLabels.value.slice(1), timeStr];
    livePressureHistory.value = [...livePressureHistory.value.slice(1), 0];
    liveFlowHistory.value = [...liveFlowHistory.value.slice(1), 0];
    liveCommandHistory.value = [...liveCommandHistory.value.slice(1), 0];
    liveFeedbackHistory.value = [...liveFeedbackHistory.value.slice(1), 0];
  }
}

// CHART DATA COMPUTED PROPERTIES
const pressureChartData = computed(() => ({
  labels: liveLabels.value, // Bind to the time array
  datasets: [
    {
      label: "Pressure (bar)",
      data: livePressureHistory.value,
      borderColor: "#3b82f6",
      backgroundColor: "rgba(59, 130, 246, 0.12)",
      fill: true,
      tension: 0.35,
      pointRadius: 0,
      borderWidth: 2,
    },
  ],
}));

const flowChartData = computed(() => ({
  labels: liveLabels.value,
  datasets: [
    {
      label: "Flow (L/min)",
      data: liveFlowHistory.value,
      borderColor: "#22c55e",
      backgroundColor: "rgba(34, 197, 94, 0.08)",
      fill: true,
      tension: 0.35,
      pointRadius: 0,
      borderWidth: 2,
    },
  ],
}));

const commandChartData = computed(() => ({
  labels: liveLabels.value,
  datasets: [
    {
      label: "Command",
      data: liveCommandHistory.value,
      borderColor: "#f59e0b",
      backgroundColor: "rgba(245, 158, 11, 0.12)",
      fill: true,
      tension: 0.35,
      pointRadius: 0,
      borderWidth: 2,
    },
  ],
}));

const feedbackChartData = computed(() => ({
  labels: liveLabels.value,
  datasets: [
    {
      label: "Feedback",
      data: liveFeedbackHistory.value,
      borderColor: "#8b5cf6",
      backgroundColor: "rgba(139, 92, 246, 0.12)",
      fill: true,
      tension: 0.35,
      pointRadius: 0,
      borderWidth: 2,
    },
  ],
}));

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: { duration: 0 },
  interaction: { intersect: false, mode: "index" },
  plugins: {
    legend: {
      display: true, 
      position: "bottom",
      labels: { 
        boxWidth: 8, 
        usePointStyle: true, 
        padding: 5, // Reduced from 12 to make the legend box tighter
        font: { size: 9 } 
      },
    },
    tooltip: { enabled: true },
  },
  scales: {
    x: {
      display: true,
      type: "category",
      grid: { display: false },
      ticks: {
        color: "#64748b",
        font: { size: 9 },
        maxTicksLimit: 8,
      },
    },
    y: {
      beginAtZero: true,
      border: { display: false },
      grid: { color: "rgba(148, 163, 184, 0.14)" },
      ticks: { color: "#64748b", font: { size: 9 } },
    },
  },
};

onMounted(async () => {
  await loadValves();
  if (route.params.id) {
    await fetchValveById(route.params.id);
  }
  window.addEventListener("keydown", handleImageViewerKeydown);
  liveDataInterval = setInterval(pollLiveData, 1000);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleImageViewerKeydown);
  document.body.style.overflow = "";
  if (liveDataInterval) clearInterval(liveDataInterval);
});

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) await fetchValveById(newId);
  },
);
</script>

<template>
  <div class="records-page">
    <div class="search-toolbar">
      <div class="search-section">
        <AutoComplete
          v-model="query"
          :suggestions="suggestions"
          optionLabel="part_number"
          placeholder="Input Part Number/Name"
          class="valve-search"
          @complete="search"
          @item-select="onSelect">
          <template #option="slotProps">
            <div class="search-option">
              <strong>{{ slotProps.option.part_number }}</strong>
              <span
                >{{ slotProps.option.manufacturer || "-" }} ·
                {{ slotProps.option.valve_type || "-" }}</span
              >
            </div>
          </template>
        </AutoComplete>
        <Button
          label="Search"
          icon="pi pi-search"
          :disabled="!pendingValve"
          @click="confirmShow" />
      </div>

      <!-- NEW: Real-time Values Display -->
      <div class="live-values-section">
                <div class="live-value-item">
          <span class="label">Real-time value</span>
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
          <span class="value"
            >{{ currentValues.pressure }} <small>bar</small></span
          >
        </div>
        <div class="live-value-item">
          <span class="label">Flow</span>
          <span class="value"
            >{{ currentValues.flow }} <small>L/min</small></span
          >
        </div>
      </div>
    </div>

    <!-- ALWAYS VISIBLE GRID LAYOUT -->
    <div class="wireframe-grid">
      <!-- LEFT COLUMN: Controls, Picture, Valve Info -->
      <div class="sidebar-column">
        <!-- Control Buttons -->
        <section class="panel control-panel">
          <Button
            label="Start"
            severity="success"
            class="ctrl-btn"
            :loading="isStarting"
            :disabled="!valve || isStarting || controlStatus === 'active'"
            @click="startValve" />
          <Button
            label="Stop"
            severity="danger"
            class="ctrl-btn"
            :loading="isStopping"
            :disabled="!valve || isStopping || controlStatus === 'stopped'"
            @click="stopValve" />
          <Button
            :label="isRecording ? 'Stop Record' : 'Record'"
            :severity="isRecording ? 'warn' : 'secondary'"
            class="ctrl-btn"
            :disabled="!valve"
            @click="toggleRecord" />
        </section>

        <!-- Picture of Valve -->
        <section class="panel picture-panel">
          <header class="panel-header">
            <h3>Picture of valve</h3>
          </header>
          <button
            v-if="imageUrl"
            type="button"
            class="picture-content picture-button"
            @click="openImageViewer">
            <img
              :src="imageUrl"
              :alt="`Valve ${valve?.part_number || ''}`"
              class="valve-image" />
          </button>
          <div v-else class="picture-content">
            <div class="image-placeholder">
              <i class="pi pi-image"></i>
              <span v-if="!valve" style="font-size: 0.7rem; margin-top: 0.5rem"
                >No valve selected</span
              >
            </div>
          </div>
        </section>

        <!-- Valve Database Info -->
        <section class="panel info-panel">
          <div class="info-row">
            <span class="info-label">Valve Name</span>
            <span class="info-val"
              >{{ valve?.manufacturer || "-" }}
              {{ valve?.valve_type || "" }}</span
            >
          </div>
          <div class="info-row">
            <span class="info-label">Part Number</span>
            <span class="info-val">{{ valve?.part_number || "-" }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Rated Flow</span>
            <span class="info-val">{{ valve?.rated_flow ?? "-" }} L/min</span>
          </div>
          <div class="info-row">
            <span class="info-label">Max Flow</span>
            <span class="info-val">{{ valve?.max_flow ?? "-" }} L/min</span>
          </div>
          <div class="info-row">
            <span class="info-label">Max Pressure</span>
            <span class="info-val">{{ valve?.max_pressure ?? "-" }} bar</span>
          </div>
        </section>
      </div>

      <!-- MIDDLE COLUMN: Pressure & Command Graphs -->
      <section class="panel graph-panel pressure-graph">
        <header class="panel-header"><h3>Pressure Graph</h3></header>
        <div class="chart-container">
          <Line
            :key="`pressure-chart`"
            :data="pressureChartData"
            :options="chartOptions" />
        </div>
      </section>

      <section class="panel graph-panel command-graph">
        <header class="panel-header"><h3>Command Graph</h3></header>
        <div class="chart-container">
          <Line
            :key="`command-chart`"
            :data="commandChartData"
            :options="chartOptions" />
        </div>
      </section>

      <!-- RIGHT COLUMN: Flow & Feedback Graphs -->
      <section class="panel graph-panel flow-graph">
        <header class="panel-header"><h3>Flow Graph</h3></header>
        <div class="chart-container">
          <Line
            :key="`flow-chart`"
            :data="flowChartData"
            :options="chartOptions" />
        </div>
      </section>

      <section class="panel graph-panel feedback-graph">
        <header class="panel-header"><h3>Feedback Graph</h3></header>
        <div class="chart-container">
          <Line
            :key="`feedback-chart`"
            :data="feedbackChartData"
            :options="chartOptions" />
        </div>
      </section>
    </div>

    <!-- Image Viewer Modal -->
    <Teleport to="body">
      <Transition name="viewer">
        <div
          v-if="imageViewerVisible"
          class="image-viewer"
          @click.self="closeImageViewer">
          <header class="viewer-header">
            <div>
              <span>VALVE IMAGE</span
              ><strong>{{ valve?.part_number || "-" }}</strong>
            </div>
            <button
              type="button"
              class="viewer-close-button"
              @click="closeImageViewer">
              <i class="pi pi-times"></i>
            </button>
          </header>
          <div class="viewer-content">
            <img
              :src="imageUrl"
              :alt="`Valve ${valve?.part_number || ''}`"
              class="viewer-image" />
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.records-page {
  width: 100%;
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  overflow: hidden;
}

.search-toolbar {
  min-height: 42px; /* Reduced from 52px */
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.3rem 0.75rem; /* Tighter internal padding */
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
  flex-wrap: wrap;
  gap: 1rem;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  width: 100%;
  max-width: 500px; /* Reduced slightly to make room for metrics */
}

/* NEW: Live Values Toolbar Styling */
.live-values-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.live-value-item {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  min-width: 4.5rem; /* Locks the width so containers never shrink/grow */
}

.live-value-item:first-child {
  min-width: auto; 
  justify-content: flex-end;
  margin-right: 0.5rem;
}

.live-value-item .label {
  font-size: 0.65rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.1rem;
}

.live-value-item .value {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--primary-color);
  font-family: Consolas, Monaco, monospace;
  font-variant-numeric: tabular-nums; /* Prevents micro-jitters between '1' and '8' widths */
  white-space: nowrap; /* Prevents the unit from wrapping to a new line */
}

.live-value-item .value small {
  font-size: 0.7rem;
  color: var(--text-muted);
  font-weight: normal;
}

.valve-search {
  flex: 1;
}
.valve-search :deep(.p-autocomplete),
.valve-search :deep(.p-inputtext) {
  width: 100%;
}

/* --- WIREFRAME GRID LAYOUT --- */
.wireframe-grid {
  flex: 1;
  min-height: 0;
  min-width: 0; /* Allows grid to shrink */
  display: grid;
  /* Use minmax(0, 1fr) instead of 1fr to prevent chart blowout */
  grid-template-columns: 280px minmax(0, 1fr) minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr) minmax(0, 1fr);
  grid-template-areas:
    "sidebar pressure flow"
    "sidebar command feedback";
  gap: 0.85rem;
}

/* Sidebar Column */
.sidebar-column {
  grid-area: sidebar;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  min-height: 0;
  min-width: 0; /* Allows column to shrink */
}

.control-panel {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.5rem;
  padding: 0.75rem;
  flex-shrink: 0;
}
.ctrl-btn {
  width: 100%;
  padding: 0.5rem 0;
}

.picture-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 200px;
}
.picture-content {
  flex: 1;
  padding: 0.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
  background: transparent;
  border: none;
}
.valve-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 4px;
}
.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--text-muted);
  font-size: 2rem;
}

.info-panel {
  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 0.75rem;
  flex-shrink: 0;
}
.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(148, 163, 184, 0.2);
  padding-bottom: 0.4rem;
}
.info-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}
.info-label {
  color: var(--text-muted);
  font-size: 0.8rem;
}
.info-val {
  color: var(--text-color);
  font-weight: 600;
  font-size: 0.9rem;
}

/* Graph Panels */
.pressure-graph {
  grid-area: pressure;
}
.flow-graph {
  grid-area: flow;
}
.command-graph {
  grid-area: command;
}
.feedback-graph {
  grid-area: feedback;
}

.graph-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  min-width: 0; /* Allows panel to shrink with grid */
}

.chart-container {
  flex: 1;
  min-height: 0;
  min-width: 0; 
  padding: 0.5rem 0.75rem 0.1rem 0.75rem; /* Reduced bottom padding from 1rem to 0.1rem */
  position: relative;
}

/* Shared Panel Styles */
.panel {
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.08);
}
.panel-header {
  padding: 0.4rem 0.8rem; /* Reduced from 0.65rem to save vertical space */
  background: rgba(148, 163, 184, 0.045);
  border-bottom: 1px solid #475569;
}
.panel-header h3 {
  margin: 0;
  color: var(--text-color);
  font-size: 0.85rem;
  font-weight: 600;
}

/* Viewer */
.image-viewer {
  position: fixed;
  z-index: 9999;
  inset: 0;
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  background: rgba(0, 0, 0, 0.9);
}
.viewer-header {
  display: flex;
  justify-content: space-between;
  color: white;
  margin-bottom: 1rem;
}
.viewer-content {
  flex: 1;
  display: flex;
  justify-content: center;
  min-height: 0;
}
.viewer-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.viewer-close-button {
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
}

:global(:root:not(.p-dark)) .panel,
:global(:root:not(.p-dark)) .search-toolbar {
  border-color: #aeb9c4;
}
:global(:root:not(.p-dark)) .panel-header {
  border-color: #c5ced7;
}

@media (max-width: 1100px) {
  .wireframe-grid {
    grid-template-columns: 250px minmax(0, 1fr);
    grid-template-rows: auto auto auto;
    grid-template-areas:
      "sidebar pressure"
      "sidebar flow"
      "sidebar command"
      "sidebar feedback";
    overflow-y: auto;
  }
}
@media (max-width: 768px) {
  .wireframe-grid {
    grid-template-columns: minmax(0, 1fr);
    grid-template-areas:
      "sidebar"
      "pressure"
      "flow"
      "command"
      "feedback";
  }
  .graph-panel {
    min-height: 250px;
  }
}
</style>
