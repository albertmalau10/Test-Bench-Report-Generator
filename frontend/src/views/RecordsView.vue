<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import api, { startOutput, stopOutput } from "../services/api";
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
  Filler,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  LinearScale,
  Filler,
);

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

async function loadValves() {
  try {
    const response = await api.get("/valves");
    allValves.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error("Gagal mengambil daftar valve:", error);
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
    console.error("Gagal mengambil detail valve:", error);
    toast.add({
      severity: "error",
      summary: "Load Failed",
      detail: "Unable to load the selected valve.",
      life: 3000,
    });
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
    const valveType = item.valve_type?.toLowerCase() || "";
    const componentSeries = item.component_series?.toLowerCase() || "";

    return (
      partNumber.includes(keyword) ||
      manufacturer.includes(keyword) ||
      valveType.includes(keyword) ||
      componentSeries.includes(keyword)
    );
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
      detail: error.response?.data?.error || "Unable to activate the output.",
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
      detail: error.response?.data?.error || "Unable to stop the output.",
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
      ? "Test data recording has started. This function is currently a placeholder."
      : "Test data recording has stopped.",
    life: 3000,
  });
}

function sendCommand() {
  if (!valve.value) return;

  toast.add({
    severity: "success",
    summary: "COMMAND",
    detail:
      `Command ${valve.value.command_type || "-"}: ` +
      `${valve.value.command_value ?? "-"} is currently a placeholder.`,
    life: 3500,
  });
}

const imageUrl = computed(() => {
  if (!valve.value?.image_path) return null
  if (valve.value.image_path.startsWith('http')) return valve.value.image_path
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${baseUrl}${valve.value.image_path}`
})

const datasheetUrl = computed(() => {
  if (!valve.value?.datasheet_path) return null
  if (valve.value.datasheet_path.startsWith('http')) return valve.value.datasheet_path
  const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${baseUrl}${valve.value.datasheet_path}`
})

function openDatasheet() {
  if (!datasheetUrl.value) return;
  window.open(datasheetUrl.value, "_blank", "noopener,noreferrer");
}

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

const STEPS = 40;

function generateSeries(target) {
  const targetValue = Number(target) || 0;
  const points = [];
  let currentValue = targetValue * 0.85;

  for (let index = 0; index < STEPS; index++) {
    const noise = (Math.random() - 0.5) * targetValue * 0.05;
    const reversion = (targetValue - currentValue) * 0.2;
    currentValue = currentValue + noise + reversion;

    points.push({
      x: index,
      y: Math.round(currentValue * 100) / 100,
    });
  }

  if (points.length > 0) {
    points[points.length - 1].y = targetValue;
  }

  return points;
}

const pressureChartData = computed(() => ({
  datasets: [
    {
      label: "Pressure",
      data: generateSeries(valve.value?.max_pressure),
      borderColor: "#3b82f6",
      backgroundColor: "rgba(59, 130, 246, 0.12)",
      fill: true,
      tension: 0.35,
      pointRadius: 0,
      pointHoverRadius: 4,
      borderWidth: 2,
    },
  ],
}));

const flowChartData = computed(() => ({
  datasets: [
    {
      label: "Rated Flow",
      data: generateSeries(valve.value?.rated_flow),
      borderColor: "#22c55e",
      backgroundColor: "rgba(34, 197, 94, 0.08)",
      fill: false,
      tension: 0.35,
      pointRadius: 0,
      pointHoverRadius: 4,
      borderWidth: 2,
    },
    {
      label: "Maximum Flow",
      data: generateSeries(valve.value?.max_flow),
      borderColor: "#3b82f6",
      backgroundColor: "rgba(59, 130, 246, 0.08)",
      fill: false,
      tension: 0.35,
      pointRadius: 0,
      pointHoverRadius: 4,
      borderWidth: 2,
    },
  ],
}));

function makeProgressiveAnimation(totalDuration = 1200) {
  const delayBetweenPoints = totalDuration / STEPS;

  const previousY = (context) => {
    if (context.index === 0) {
      return context.chart.scales.y.getPixelForValue(0);
    }

    const previousPoint = context.chart.getDatasetMeta(context.datasetIndex)
      .data[context.index - 1];

    return previousPoint?.getProps(["y"], true).y;
  };

  return {
    x: {
      type: "number",
      easing: "easeOutQuad",
      duration: delayBetweenPoints,
      from: Number.NaN,
      delay(context) {
        if (context.type !== "data" || context.xStarted) return 0;
        context.xStarted = true;
        return context.index * delayBetweenPoints;
      },
    },
    y: {
      type: "number",
      easing: "easeOutQuad",
      duration: delayBetweenPoints,
      from: previousY,
      delay(context) {
        if (context.type !== "data" || context.yStarted) return 0;
        context.yStarted = true;
        return context.index * delayBetweenPoints;
      },
    },
  };
}

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: makeProgressiveAnimation(),
  interaction: {
    intersect: false,
    mode: "index",
  },
  plugins: {
    legend: {
      display: true,
      position: "bottom",
      labels: {
        boxWidth: 10,
        boxHeight: 3,
        usePointStyle: true,
        pointStyle: "line",
        padding: 12,
        font: { size: 10 },
      },
    },
    tooltip: { enabled: true },
  },
  scales: {
    x: {
      display: false,
      type: "linear",
      grid: { display: false },
    },
    y: {
      beginAtZero: true,
      border: { display: false },
      grid: { color: "rgba(148, 163, 184, 0.14)" },
      ticks: {
        color: "#64748b",
        font: { size: 9 },
      },
    },
  },
};

const feedback = ref({
  position: 65,
  current: 1.25,
  temperature: 42,
  status: "Normal",
});

onMounted(async () => {
  await loadValves();

  if (route.params.id) {
    await fetchValveById(route.params.id);
  }

  window.addEventListener("keydown", handleImageViewerKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleImageViewerKeydown);
  document.body.style.overflow = "";
});

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await fetchValveById(newId);
    }
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
          placeholder="Search valve by part number"
          class="valve-search"
          @complete="search"
          @item-select="onSelect">
          <template #option="slotProps">
            <div class="search-option">
              <strong>{{ slotProps.option.part_number }}</strong>
              <span>
                {{ slotProps.option.manufacturer || "-" }}
                ·
                {{ slotProps.option.valve_type || "-" }}
              </span>
            </div>
          </template>
        </AutoComplete>

        <Button
          label="Show Valve"
          icon="pi pi-search"
          :disabled="!pendingValve"
          @click="confirmShow" />
      </div>

      <div v-if="valve" class="selected-valve">
        <span>SELECTED VALVE</span>
        <strong>{{ valve.part_number }}</strong>
      </div>
    </div>

    <div v-if="!valve" class="empty-state">
      <i class="pi pi-sliders-h"></i>
      <h3>No valve selected</h3>
      <p>
        Search for a valve by part number and select Show Valve to open the test
        workspace.
      </p>
    </div>

    <div v-else class="records-grid">
      <section class="panel picture-panel">
        <header class="panel-header">
          <h3>Picture of Valve</h3>

          <button
            v-if="imageUrl"
            type="button"
            class="panel-action"
            title="Open image viewer"
            aria-label="Open valve image viewer"
            @click="openImageViewer">
            <i class="pi pi-window-maximize"></i>
          </button>
        </header>

        <button
          v-if="imageUrl"
          type="button"
          class="picture-content picture-button"
          title="Click to enlarge"
          @click="openImageViewer">
          <img
            :src="imageUrl"
            :alt="`Valve ${valve.part_number}`"
            class="valve-image" />

          <span class="image-expand-hint">
            <i class="pi pi-search-plus"></i>
            Enlarge
          </span>
        </button>

        <div v-else class="picture-content">
          <div class="image-placeholder">
            <i class="pi pi-image"></i>
            <span>No image available</span>
          </div>
        </div>

        <footer class="valve-information">
          <strong>{{ valve.part_number }}</strong>
          <span>
            {{ valve.component_series || "-" }}
            ·
            {{ valve.valve_type || "-" }}
          </span>
        </footer>
      </section>

      <section class="panel datasheet-panel">
        <header class="panel-header">
          <h3>Valve Datasheet</h3>
          <i class="pi pi-file-pdf panel-icon pdf-icon"></i>
        </header>

        <div class="datasheet-content">
          <template v-if="datasheetUrl">
            <div class="document-symbol">
              <i class="pi pi-file-pdf"></i>
            </div>

            <div class="document-information">
              <strong>Technical Datasheet</strong>
              <span>{{ valve.part_number }}.pdf</span>
            </div>

            <Button
              label="View PDF"
              icon="pi pi-external-link"
              severity="secondary"
              outlined
              class="datasheet-button"
              @click="openDatasheet" />
          </template>

          <template v-else>
            <div class="document-symbol unavailable">
              <i class="pi pi-file"></i>
            </div>

            <div class="document-information">
              <strong>Not Available</strong>
              <span>No datasheet has been uploaded.</span>
            </div>
          </template>
        </div>
      </section>

      <section class="panel pressure-panel">
        <header class="panel-header">
          <h3>Pressure</h3>
          <div class="primary-metric">
            <strong>{{ valve.max_pressure ?? "-" }}</strong>
            <span>bar</span>
          </div>
        </header>

        <div class="chart-container">
          <Line
            :key="`${valve.id}-pressure`"
            :data="pressureChartData"
            :options="chartOptions" />
        </div>
      </section>

      <section class="panel flow-panel">
        <header class="panel-header">
          <h3>Flow</h3>
          <div class="flow-metrics">
            <div>
              <span>Rated</span>
              <strong>{{ valve.rated_flow ?? "-" }}</strong>
            </div>
            <div>
              <span>Maximum</span>
              <strong>{{ valve.max_flow ?? "-" }}</strong>
            </div>
          </div>
        </header>

        <div class="chart-container">
          <Line
            :key="`${valve.id}-flow`"
            :data="flowChartData"
            :options="chartOptions" />
        </div>
      </section>

      <section class="panel command-panel">
        <header class="panel-header">
          <h3>Command</h3>
          <i class="pi pi-send panel-icon"></i>
        </header>

        <div class="command-content">
          <div class="command-item">
            <span>Command Type</span>
            <strong>{{ valve.command_type || "-" }}</strong>
          </div>
          <div class="command-divider"></div>
          <div class="command-item">
            <span>Command Value</span>
            <strong>{{ valve.command_value ?? "-" }}</strong>
          </div>
        </div>
      </section>

      <section class="panel feedback-panel">
        <header class="panel-header">
          <h3>Feedback</h3>
          <div class="feedback-condition">
            <span class="feedback-dot"></span>
            {{ feedback.status }}
          </div>
        </header>

        <div class="feedback-content">
          <div class="feedback-item">
            <span>Position</span>
            <div>
              <strong>{{ feedback.position }}</strong>
              <small>%</small>
            </div>
          </div>
          <div class="feedback-item">
            <span>Current</span>
            <div>
              <strong>{{ feedback.current }}</strong>
              <small>A</small>
            </div>
          </div>
          <div class="feedback-item">
            <span>Temperature</span>
            <div>
              <strong>{{ feedback.temperature }}</strong>
              <small>°C</small>
            </div>
          </div>
        </div>

        <footer class="simulation-warning">
          Simulation data, not connected in real time.
        </footer>
      </section>

      <section class="control-panel">
        <div class="control-buttons">
          <Button
            label="Start"
            icon="pi pi-play"
            severity="success"
            class="control-button"
            :loading="isStarting"
            :disabled="!valve || isStarting || controlStatus === 'active'"
            @click="startValve" />

          <Button
            label="Stop"
            icon="pi pi-stop"
            severity="danger"
            outlined
            class="control-button"
            :loading="isStopping"
            :disabled="!valve || isStopping || controlStatus === 'stopped'"
            @click="stopValve" />

          <Button
            :label="isRecording ? 'Stop Record' : 'Record'"
            :icon="isRecording ? 'pi pi-stop-circle' : 'pi pi-circle-fill'"
            :severity="isRecording ? 'warn' : 'secondary'"
            outlined
            class="control-button"
            :disabled="!valve"
            @click="toggleRecord" />

          <Button
            label="Send Command"
            icon="pi pi-send"
            class="control-button"
            :disabled="!valve"
            @click="sendCommand" />
        </div>

        <div class="output-status">
          <span class="output-label">OUTPUT STATUS</span>
          <div class="status-badge" :class="controlStatus">
            <span class="status-dot"></span>
            <span>{{ controlStatus === "active" ? "Active" : "Stopped" }}</span>
          </div>
        </div>
      </section>
    </div>

    <Teleport to="body">
      <Transition name="viewer">
        <div
          v-if="imageViewerVisible"
          class="image-viewer"
          role="dialog"
          aria-modal="true"
          aria-label="Valve image viewer"
          @click.self="closeImageViewer">
          <header class="viewer-header">
            <div>
              <span>VALVE IMAGE</span>
              <strong>{{ valve?.part_number || "-" }}</strong>
            </div>

            <button
              type="button"
              class="viewer-close-button"
              aria-label="Close image viewer"
              title="Close image viewer"
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

          <footer class="viewer-footer">
            <span>
              {{ valve?.component_series || "-" }}
              ·
              {{ valve?.valve_type || "-" }}
            </span>
            <span>Press ESC or click outside the image to close</span>
          </footer>
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
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex: 0 0 auto;
  padding: 0.5rem;
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
}

.search-section {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 0.65rem;
  flex: 1;
}

.valve-search {
  min-width: 0;
  max-width: 36rem;
  flex: 1;
}

.valve-search :deep(.p-autocomplete),
.valve-search :deep(.p-inputtext) {
  width: 100%;
}

.search-option {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  padding: 0.2rem 0;
}

.search-option strong {
  color: var(--text-color);
  font-size: 0.82rem;
}

.search-option span {
  color: var(--text-muted);
  font-size: 0.7rem;
}

.selected-valve {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.1rem;
  padding: 0 0.5rem;
}

.selected-valve span,
.output-label {
  color: var(--text-muted);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.55rem;
  letter-spacing: 0.08em;
}

.selected-valve strong {
  color: var(--primary-color);
  font-size: 0.78rem;
}

.records-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns:
    minmax(145px, 0.58fr)
    minmax(260px, 1fr)
    minmax(290px, 1.12fr);
  grid-template-rows:
    minmax(180px, 1fr)
    minmax(165px, 0.95fr)
    auto;
  grid-template-areas:
    "picture pressure flow"
    "datasheet command feedback"
    "datasheet controls controls";
  gap: 0.85rem;
}

.picture-panel {
  grid-area: picture;
}
.datasheet-panel {
  grid-area: datasheet;
}
.pressure-panel {
  grid-area: pressure;
}
.flow-panel {
  grid-area: flow;
}
.command-panel {
  grid-area: command;
}
.feedback-panel {
  grid-area: feedback;
}
.control-panel {
  grid-area: controls;
}

.panel {
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.08);
}

.panel-header {
  min-height: 46px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex: 0 0 auto;
  padding: 0.65rem 0.8rem;
  background: rgba(148, 163, 184, 0.045);
  border-bottom: 1px solid #475569;
}

.panel h3 {
  margin: 0;
  color: var(--text-color);
  font-size: 0.8rem;
  font-weight: 650;
  letter-spacing: 0.01em;
}

.panel-icon {
  flex-shrink: 0;
  color: var(--text-muted);
  font-size: 0.88rem;
}

.pdf-icon {
  color: #dc2626;
}

.panel-action {
  width: 1.85rem;
  height: 1.85rem;
  display: grid;
  place-items: center;
  padding: 0;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid transparent;
  border-radius: 3px;
  cursor: pointer;
}

.panel-action:hover {
  color: var(--text-color);
  background: var(--sidebar-active-bg);
  border-color: var(--border-color);
}

.panel-action:focus-visible,
.picture-button:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
}

.picture-content {
  position: relative;
  flex: 1;
  min-height: 0;
  padding: 0.55rem;
}

.picture-button {
  width: 100%;
  display: block;
  appearance: none;
  color: inherit;
  background: transparent;
  border: none;
  cursor: zoom-in;
  text-align: inherit;
}

.valve-image {
  width: 100%;
  height: 100%;
  min-height: 80px;
  display: block;
  object-fit: contain;
  background: var(--bg-color);
  border: 1px solid var(--border-color);
  border-radius: 2px;
}

.image-expand-hint {
  position: absolute;
  right: 0.9rem;
  bottom: 0.9rem;
  display: inline-flex;
  align-items: center;
  gap: 0.32rem;
  padding: 0.28rem 0.45rem;
  color: #e2e8f0;
  background: rgba(15, 23, 42, 0.82);
  border: 1px solid rgba(148, 163, 184, 0.4);
  border-radius: 2px;
  font-size: 0.58rem;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.15s ease;
}

.picture-button:hover .image-expand-hint,
.picture-button:focus-visible .image-expand-hint {
  opacity: 1;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  min-height: 90px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  color: var(--text-muted);
  background: var(--bg-color);
  border: 1px dashed #64748b;
  border-radius: 2px;
}

.image-placeholder i {
  font-size: 1.45rem;
}
.image-placeholder span {
  font-size: 0.65rem;
}

.valve-information {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  flex: 0 0 auto;
  padding: 0.65rem 0.8rem;
  border-top: 1px solid #475569;
}

.valve-information strong {
  overflow: hidden;
  color: var(--primary-color);
  font-size: 0.8rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.valve-information span {
  overflow: hidden;
  color: var(--text-muted);
  font-size: 0.62rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.datasheet-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.8rem;
  padding: 1rem;
  text-align: center;
}

.document-symbol {
  width: 3.2rem;
  height: 3.2rem;
  display: grid;
  place-items: center;
  color: #dc2626;
  background: rgba(220, 38, 38, 0.09);
  border: 1px solid rgba(220, 38, 38, 0.28);
  border-radius: 3px;
}

.document-symbol i {
  font-size: 1.45rem;
}

.document-symbol.unavailable {
  color: var(--text-muted);
  background: var(--bg-color);
  border-color: #475569;
}

.document-information {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.document-information strong {
  color: var(--text-color);
  font-size: 0.76rem;
}

.document-information span {
  overflow: hidden;
  color: var(--text-muted);
  font-size: 0.64rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.datasheet-button {
  width: 100%;
  max-width: 9rem;
}

.chart-container {
  position: relative;
  flex: 1;
  min-height: 0;
  padding: 0.55rem 0.7rem 0.65rem;
}

.primary-metric {
  display: flex;
  align-items: baseline;
  gap: 0.3rem;
}

.primary-metric strong {
  color: var(--text-color);
  font-size: 1.1rem;
}

.primary-metric span {
  color: var(--text-muted);
  font-size: 0.62rem;
}

.flow-metrics {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.flow-metrics > div {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.05rem;
}

.flow-metrics span {
  color: var(--text-muted);
  font-size: 0.52rem;
  text-transform: uppercase;
}

.flow-metrics strong {
  color: var(--text-color);
  font-size: 0.8rem;
}

.command-content {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 1rem;
  padding: 1.15rem;
}

.command-item {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.command-item span,
.feedback-item > span {
  color: var(--text-muted);
  font-size: 0.62rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.command-item strong {
  overflow: hidden;
  color: var(--text-color);
  font-family: Consolas, Monaco, monospace;
  font-size: 1.15rem;
  text-overflow: ellipsis;
}

.command-divider {
  width: 1px;
  height: 3rem;
  background: #475569;
}

.feedback-condition {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: #16a34a;
  font-size: 0.66rem;
  font-weight: 600;
}

.feedback-dot,
.status-dot {
  width: 7px;
  height: 7px;
  display: inline-block;
  flex: 0 0 auto;
  background: currentColor;
  border-radius: 50%;
}

.feedback-content {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  align-items: center;
  padding: 0.9rem;
}

.feedback-item {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  padding: 0 0.7rem;
  border-right: 1px solid #475569;
}

.feedback-item:first-child {
  padding-left: 0;
}
.feedback-item:last-child {
  padding-right: 0;
  border-right: none;
}

.feedback-item > div {
  display: flex;
  align-items: baseline;
  gap: 0.2rem;
}

.feedback-item strong {
  color: var(--text-color);
  font-family: Consolas, Monaco, monospace;
  font-size: 1.1rem;
}

.feedback-item small {
  color: var(--text-muted);
  font-size: 0.62rem;
}

.simulation-warning {
  flex: 0 0 auto;
  padding: 0.45rem 0.9rem;
  color: var(--text-muted);
  border-top: 1px solid #475569;
  font-size: 0.57rem;
}

.control-panel {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.7rem;
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.08);
}

.control-buttons {
  min-width: 0;
  display: grid;
  grid-template-columns:
    minmax(88px, 0.7fr)
    minmax(88px, 0.7fr)
    minmax(110px, 0.85fr)
    minmax(155px, 1.7fr);
  gap: 0.6rem;
  flex: 1;
}

.control-button {
  width: 100%;
}

.output-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.28rem;
  flex: 0 0 auto;
}

.status-badge {
  min-width: 6.3rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  padding: 0.35rem 0.65rem;
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 650;
}

.status-badge.stopped {
  color: #dc2626;
  background: rgba(220, 38, 38, 0.09);
  border: 1px solid rgba(220, 38, 38, 0.3);
}

.status-badge.active {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.09);
  border: 1px solid rgba(22, 163, 74, 0.3);
}

.empty-state {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: var(--text-muted);
  background: var(--card-bg);
  border: 1px solid #475569;
  border-radius: 4px;
  text-align: center;
}

.empty-state > i {
  margin-bottom: 1rem;
  color: var(--primary-color);
  font-size: 1.5rem;
}

.empty-state h3 {
  margin: 0 0 0.45rem;
  color: var(--text-color);
  font-size: 1rem;
}

.empty-state p {
  max-width: 25rem;
  margin: 0;
  font-size: 0.78rem;
  line-height: 1.6;
}

.image-viewer {
  position: fixed;
  z-index: 9999;
  inset: 0;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  padding: 1.25rem;
  color: #e5e7eb;
  background: rgba(5, 8, 11, 0.96);
  backdrop-filter: blur(4px);
}

.viewer-header {
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding-bottom: 0.85rem;
  border-bottom: 1px solid #4b5563;
}

.viewer-header > div {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.18rem;
}

.viewer-header span {
  color: #7b8792;
  font-family: Consolas, Monaco, monospace;
  font-size: 0.58rem;
  letter-spacing: 0.1em;
}

.viewer-header strong {
  overflow: hidden;
  color: #f1f5f9;
  font-size: 0.88rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.viewer-close-button {
  width: 2.35rem;
  height: 2.35rem;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  padding: 0;
  color: #cbd5e1;
  background: transparent;
  border: 1px solid #64748b;
  border-radius: 3px;
  cursor: pointer;
}

.viewer-close-button:hover {
  color: #ffffff;
  background: #252c33;
}

.viewer-content {
  min-width: 0;
  min-height: 0;
  display: grid;
  place-items: center;
  padding: 1.5rem;
  overflow: hidden;
}

.viewer-image {
  display: block;
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  background: #11161a;
  border: 1px solid #64748b;
}

.viewer-footer {
  min-height: 42px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding-top: 0.8rem;
  color: #76838e;
  border-top: 1px solid #4b5563;
  font-size: 0.64rem;
}

.viewer-enter-active,
.viewer-leave-active {
  transition: opacity 0.15s ease;
}

.viewer-enter-from,
.viewer-leave-to {
  opacity: 0;
}

:global(:root:not(.p-dark)) .panel,
:global(:root:not(.p-dark)) .control-panel,
:global(:root:not(.p-dark)) .search-toolbar,
:global(:root:not(.p-dark)) .empty-state {
  border-color: #aeb9c4;
}

:global(:root:not(.p-dark)) .panel-header,
:global(:root:not(.p-dark)) .valve-information,
:global(:root:not(.p-dark)) .simulation-warning,
:global(:root:not(.p-dark)) .feedback-item,
:global(:root:not(.p-dark)) .command-divider {
  border-color: #c5ced7;
}

@media (max-width: 1100px) {
  .records-page {
    overflow-y: auto;
  }

  .records-grid {
    flex: none;
    min-height: 560px;
    grid-template-columns:
      minmax(145px, 0.62fr)
      minmax(230px, 1fr)
      minmax(250px, 1fr);
    grid-template-rows: 220px 190px auto;
  }

  .output-status {
    display: none;
  }
}

@media (max-width: 820px) {
  .records-grid {
    min-height: auto;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    grid-template-rows: auto;
    grid-template-areas:
      "picture datasheet"
      "pressure pressure"
      "flow flow"
      "command feedback"
      "controls controls";
  }

  .picture-panel,
  .datasheet-panel {
    min-height: 230px;
  }

  .pressure-panel,
  .flow-panel {
    min-height: 260px;
  }

  .command-panel,
  .feedback-panel {
    min-height: 190px;
  }

  .control-buttons {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 580px) {
  .records-page {
    overflow-y: auto;
  }

  .search-toolbar,
  .search-section,
  .control-panel {
    align-items: stretch;
    flex-direction: column;
  }

  .valve-search {
    width: 100%;
    max-width: none;
  }

  .selected-valve {
    align-items: flex-start;
  }

  .records-grid {
    display: flex;
    min-height: auto;
    flex-direction: column;
  }

  .picture-panel,
  .datasheet-panel,
  .command-panel,
  .feedback-panel {
    min-height: 215px;
  }

  .pressure-panel,
  .flow-panel {
    min-height: 280px;
  }

  .control-buttons {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .output-status {
    display: flex;
    align-items: stretch;
  }

  .status-badge {
    width: 100%;
  }
  .feedback-item {
    padding: 0 0.4rem;
  }
  .flow-metrics {
    gap: 0.55rem;
  }
  .image-expand-hint {
    opacity: 1;
  }

  .image-viewer {
    padding: 0.75rem;
  }
  .viewer-content {
    padding: 0.75rem 0;
  }

  .viewer-footer {
    align-items: flex-start;
    flex-direction: column;
    justify-content: center;
    gap: 0.2rem;
  }
}

@media (max-width: 390px) {
  .control-buttons {
    grid-template-columns: 1fr;
  }

  .feedback-content {
    grid-template-columns: 1fr;
    gap: 0.7rem;
  }

  .feedback-item {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 0 0 0.55rem;
    border-right: none;
    border-bottom: 1px solid #475569;
  }

  .feedback-item:last-child {
    padding-bottom: 0;
    border-bottom: none;
  }
}
</style>
