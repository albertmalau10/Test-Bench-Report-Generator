<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import api from "../services/api";
import InputText from "primevue/inputtext";
import Dropdown from "primevue/dropdown";
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

const toast = useToast();

// Form & History State
const historicalRecords = ref([]);
const selectedRecord = ref(null);
const docNumber = ref("");
const operatorName = ref("");
const selectedStatus = ref("OK");
const statusOptions = [
  { label: "OK (Passed)", value: "OK" },
  { label: "NOK (Failed)", value: "NOK" },
  { label: "Requires Calibration", value: "CALIBRATION" },
];

// Computed Data Extraction
const valve = computed(() => selectedRecord.value?.Valve || null);

const reportDate = computed(() => {
  if (!selectedRecord.value) return "-";
  return new Date(selectedRecord.value.created_at).toLocaleString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
});

const recordOptions = computed(() => {
  return historicalRecords.value.map((r) => ({
    label: `${r.Valve?.part_number || "Unknown"} — ${new Date(r.created_at).toLocaleString("en-GB")}`,
    value: r,
  }));
});

// Parse the saved telemetry JSON
const parsedTelemetry = computed(() => {
  if (!selectedRecord.value?.data_payload) return [];
  try {
    return JSON.parse(selectedRecord.value.data_payload);
  } catch (e) {
    return [];
  }
});

// Dynamic Parameter Verification
const maxRecordedPressure = computed(() => {
  if (!parsedTelemetry.value.length) return 0;
  return Math.max(
    ...parsedTelemetry.value.map((d) => Number(d.pressure)),
  ).toFixed(1);
});

const maxRecordedFlow = computed(() => {
  if (!parsedTelemetry.value.length) return 0;
  return Math.max(...parsedTelemetry.value.map((d) => Number(d.flow))).toFixed(
    1,
  );
});

const pressureVerification = computed(() => {
  if (!valve.value) return "-";
  return Number(maxRecordedPressure.value) <= valve.value.max_pressure
    ? "Within Tolerance"
    : "Exceeded Limits";
});

const flowVerification = computed(() => {
  if (!valve.value) return "-";
  return Number(maxRecordedFlow.value) <= valve.value.max_flow
    ? "Within Tolerance"
    : "Exceeded Limits";
});

// Chart Bindings
const hydraulicChartData = computed(() => {
  const data = parsedTelemetry.value;
  return {
    labels: data.map((d) => d.time),
    datasets: [
      {
        label: "Pressure (bar)",
        data: data.map((d) => d.pressure),
        borderColor: "#00CCFF",
        backgroundColor: "rgba(0, 204, 255, 0.1)",
        fill: true,
        tension: 0.35,
        pointRadius: 0,
        borderWidth: 2,
      },
      {
        label: "Flow (L/min)",
        data: data.map((d) => d.flow),
        borderColor: "#10b981",
        backgroundColor: "rgba(16, 185, 129, 0.1)",
        fill: true,
        tension: 0.35,
        pointRadius: 0,
        borderWidth: 2,
      },
    ],
  };
});

const electricalChartData = computed(() => {
  const data = parsedTelemetry.value;
  return {
    labels: data.map((d) => d.time),
    datasets: [
      {
        label: "Command",
        data: data.map((d) => d.command),
        borderColor: "#f59e0b",
        backgroundColor: "rgba(245, 158, 11, 0.1)",
        fill: true,
        tension: 0.35,
        pointRadius: 0,
        borderWidth: 2,
      },
      {
        label: "Feedback",
        data: data.map((d) => d.feedback),
        borderColor: "#8b5cf6",
        backgroundColor: "rgba(139, 92, 246, 0.1)",
        fill: true,
        tension: 0.35,
        pointRadius: 0,
        borderWidth: 2,
      },
    ],
  };
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  plugins: {
    legend: { position: "bottom", labels: { boxWidth: 10, font: { size: 9 } } },
    tooltip: { enabled: false }, 
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { maxTicksLimit: 6, font: { size: 8 } },
    },
    y: {
      beginAtZero: true,
      grid: { color: "rgba(0,0,0,0.05)" },
      ticks: { font: { size: 8 } },
    },
  },
};

async function loadHistory() {
  try {
    const res = await api.get("/records");
    historicalRecords.value = res.data;
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: "Failed to load test history.",
      life: 3000,
    });
  }
}

async function generatePDF() {
  if (!selectedRecord.value) {
    toast.add({
      severity: "error",
      summary: "Missing Data",
      detail: "Please select a test record from the history.",
      life: 3000,
    });
    return;
  }

  if (!docNumber.value.trim()) {
    docNumber.value = `DOC-${Date.now().toString().slice(-6)}`;
  }

  await nextTick();
  window.print();
}

onMounted(() => {
  loadHistory();
});
</script>

<template>
  <div class="report-generator-page">
    <!-- LEFT: Controls (Hidden on Print) -->
    <aside class="control-panel no-print">
      <header class="panel-header">
        <h2>Report Configuration</h2>
        <p>Load history and generate PDF.</p>
      </header>

      <div class="form-grid">
        <div class="field">
          <label>Select Test Record (History)</label>
          <Dropdown
            v-model="selectedRecord"
            :options="recordOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="Select historical data"
            class="w-full"
            filter />
        </div>

        <div class="field">
          <label>Document Number (Optional)</label>
          <InputText
            v-model="docNumber"
            placeholder="Leave empty to auto-generate" />
        </div>

        <div class="field">
          <label>Operator Name</label>
          <InputText v-model="operatorName" placeholder="Enter operator name" />
        </div>

        <div class="field">
          <label>Final Test Status</label>
          <Dropdown
            v-model="selectedStatus"
            :options="statusOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full" />
        </div>

        <div v-if="parsedTelemetry.length > 0" class="valve-summary-card">
          <strong>{{ parsedTelemetry.length }} Data Points Loaded</strong>
          <span>Ready for rendering.</span>
        </div>
      </div>

      <div class="action-footer">
        <Button
          label="Generate PDF Report"
          icon="pi pi-file-pdf"
          size="large"
          class="w-full"
          @click="generatePDF" />
      </div>
    </aside>

    <!-- RIGHT: A4 Document Preview -->
    <main class="document-viewer">
      <div class="a4-paper" id="print-area">
        <!-- Report Header -->
        <header class="report-header">
          <div class="brand-block">
            <h1 class="brand-title">Bosch Rexroth</h1>
            <span class="brand-subtitle">Industrial Hydraulics</span>
          </div>
          <div class="doc-title">
            <h2>VALVE TEST REPORT</h2>
            <span class="doc-id">{{ docNumber || "PENDING" }}</span>
          </div>
        </header>

        <hr class="divider" />

        <!-- Meta Information -->
        <section class="info-section">
          <div class="info-column">
            <div class="info-row">
              <span class="label">Date / Time:</span>
              <span class="value">{{ reportDate }}</span>
            </div>
            <div class="info-row">
              <span class="label">Operator:</span>
              <span class="value">{{ operatorName || "-" }}</span>
            </div>
            <div class="info-row">
              <span class="label">Test System:</span>
              <span class="value">ctrlX CORE Test Bench</span>
            </div>
          </div>
          <div class="info-column status-column">
            <div class="info-row">
              <span class="label">Overall Status:</span>
              <span class="status-badge" :class="selectedStatus">{{
                selectedStatus
              }}</span>
            </div>
          </div>
        </section>

        <!-- Valve Specifications Table -->
        <section class="spec-section">
          <h3>1. Component Specifications</h3>
          <table class="spec-table" v-if="valve">
            <tbody>
              <tr>
                <th width="20%">Manufacturer</th>
                <td width="30%">{{ valve.manufacturer }}</td>
                <th width="20%">Part Number</th>
                <td width="30%">
                  <strong>{{ valve.part_number }}</strong>
                </td>
              </tr>
              <tr>
                <th>Series</th>
                <td>{{ valve.component_series }}</td>
                <th>Valve Type</th>
                <td>{{ valve.valve_type }}</td>
              </tr>
              <tr>
                <th>Size NG</th>
                <td>{{ valve.size_ng }}</td>
                <th>Unit Weight</th>
                <td>{{ valve.weight }} kg</td>
              </tr>
            </tbody>
          </table>
          <div v-else class="empty-state">No historical record selected.</div>
        </section>

        <!-- Reference vs Tested Table -->
        <section class="spec-section" v-if="valve">
          <h3>2. Parameter Verification</h3>
          <table class="spec-table verification-table">
            <thead>
              <tr>
                <th>Parameter</th>
                <th>Database Reference</th>
                <th>Peak Recorded Value</th>
                <th>Deviation Status</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Max Pressure</td>
                <td>{{ valve.max_pressure }} bar</td>
                <td>{{ maxRecordedPressure }} bar</td>
                <td
                  :class="
                    pressureVerification === 'Within Tolerance'
                      ? 'text-ok'
                      : 'text-fail'
                  ">
                  {{ pressureVerification }}
                </td>
              </tr>
              <tr>
                <td>Max Flow</td>
                <td>{{ valve.max_flow }} L/min</td>
                <td>{{ maxRecordedFlow }} L/min</td>
                <td
                  :class="
                    flowVerification === 'Within Tolerance'
                      ? 'text-ok'
                      : 'text-fail'
                  ">
                  {{ flowVerification }}
                </td>
              </tr>
              <tr>
                <td>Command Input Type</td>
                <td>{{ valve.command_type || "Voltage" }}</td>
                <td>{{ valve.command_type || "Voltage" }}</td>
                <td class="text-ok">Matched</td>
              </tr>
            </tbody>
          </table>
        </section>

        <!-- Performance Graphs -->
        <section class="graphs-section" v-if="parsedTelemetry.length > 0">
          <h3>3. Dynamic Test Telemetry</h3>
          <div class="graph-row">
            <div class="graph-box">
              <h4>Hydraulic Performance (Pressure & Flow)</h4>
              <div class="chart-wrapper">
                <Line :data="hydraulicChartData" :options="chartOptions" />
              </div>
            </div>
          </div>
          <div class="graph-row mt-4">
            <div class="graph-box">
              <h4>Electrical Response (Command vs Feedback)</h4>
              <div class="chart-wrapper">
                <Line :data="electricalChartData" :options="chartOptions" />
              </div>
            </div>
          </div>
        </section>
        <div v-else class="empty-state">
          Graphs will render automatically when a historical test session is
          loaded.
        </div>

        <!-- Signatures -->
        <section class="signature-section">
          <div class="sig-box">
            <span>Tested By (Operator)</span>
            <div class="sig-line"></div>
            <span class="sig-name">{{ operatorName }}</span>
          </div>
          <div class="sig-box">
            <span>Approved By (Supervisor)</span>
            <div class="sig-line"></div>
            <span class="sig-name"></span>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<style scoped>
.report-generator-page {
  display: flex;
  height: 100%;
  gap: 1.5rem;
  overflow: hidden;
}

/* Control Panel */
.control-panel {
  width: 360px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel-header {
  padding: 1.5rem 1.5rem 1rem;
  border-bottom: 1px solid var(--border-color);
}
.panel-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-color);
}
.panel-header p {
  margin: 0.25rem 0 0;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.form-grid {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  overflow-y: auto;
  flex: 1;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}
.field label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
}

.valve-summary-card {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid #10b981;
  padding: 0.75rem;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}
.valve-summary-card strong {
  color: #10b981;
  font-size: 0.9rem;
}
.valve-summary-card span {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.action-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--border-color);
}
.w-full {
  width: 100%;
}
.mt-4 {
  margin-top: 1rem;
}
.empty-state {
  padding: 2rem;
  text-align: center;
  border: 1px dashed var(--border-color);
  color: var(--text-muted);
  font-size: 0.9rem;
}

/* Document Viewer / A4 Paper */
.document-viewer {
  flex: 1;
  overflow-y: auto;
  display: flex;
  justify-content: center;
  padding: 1rem;
  background: rgba(0, 0, 0, 0.02);
}

.a4-paper {
  width: 210mm;
  min-height: 297mm;
  background: white;
  padding: 15mm 20mm;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  color: #000;
  font-family: "Segoe UI", Arial, sans-serif;
  margin: 0 auto;
}

/* Report Internal Styling (Strictly Black/Dark Blue for printing) */
.report-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}
.brand-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 900;
  color: #002b49;
  letter-spacing: -1px;
}
.brand-subtitle {
  font-size: 0.9rem;
  color: #43545f;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.doc-title {
  text-align: right;
}
.doc-title h2 {
  margin: 0;
  font-size: 1.4rem;
  color: #002b49;
}
.doc-id {
  font-family: Consolas, monospace;
  font-size: 0.85rem;
  color: #666;
  font-weight: bold;
}

.divider {
  border: none;
  border-top: 2px solid #002b49;
  margin: 0 0 1.5rem 0;
}

.info-section {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
}
.info-column {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}
.info-row {
  display: flex;
  gap: 0.5rem;
}
.info-row .label {
  width: 100px;
  font-weight: 600;
  color: #43545f;
}
.info-row .value {
  font-weight: 500;
}
.status-column {
  align-items: flex-end;
}

.status-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
  font-weight: 700;
  font-size: 1.1rem;
  border: 2px solid #000;
}
.status-badge.OK {
  color: #10b981;
  border-color: #10b981;
}
.status-badge.NOK {
  color: #ef4444;
  border-color: #ef4444;
}
.status-badge.CALIBRATION {
  color: #f59e0b;
  border-color: #f59e0b;
}

h3 {
  font-size: 1.1rem;
  color: #002b49;
  border-bottom: 1px solid #ccc;
  padding-bottom: 0.25rem;
  margin: 0 0 0.75rem 0;
}

.spec-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1.5rem;
  font-size: 0.85rem;
}
.spec-table th,
.spec-table td {
  border: 1px solid #ddd;
  padding: 0.4rem 0.6rem;
  text-align: left;
}
.spec-table th {
  background: #f8fafc;
  font-weight: 600;
  color: #43545f;
}
.text-ok {
  color: #10b981;
  font-weight: 600;
}
.text-fail {
  color: #ef4444;
  font-weight: 600;
}

.graphs-section {
  margin-bottom: 2rem;
}
.graph-box h4 {
  margin: 0 0 0.5rem 0;
  font-size: 0.85rem;
  color: #43545f;
}
.chart-wrapper {
  height: 200px;
  width: 100%;
  border: 1px solid #eee;
  padding: 0.5rem;
}

.signature-section {
  display: flex;
  justify-content: space-between;
  margin-top: 3rem;
}
.sig-box {
  width: 40%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
}
.sig-line {
  width: 100%;
  border-bottom: 1px solid #000;
  margin-top: 2rem;
}
.sig-name {
  font-weight: 600;
  min-height: 1.2rem;
}

/* PRINT MEDIA QUERIES */
@media print {
  @page {
    size: A4 portrait;
    margin: 0;
  }
  body * {
    visibility: hidden;
  }
  .no-print {
    display: none !important;
  }

  .document-viewer {
    position: absolute;
    left: 0;
    top: 0;
    width: 210mm;
    height: 297mm;
    padding: 0;
    background: white;
    overflow: visible;
  }

  .a4-paper,
  .a4-paper * {
    visibility: visible;
  }

  .a4-paper {
    position: absolute;
    left: 0;
    top: 0;
    box-shadow: none;
    padding: 15mm 20mm;
    width: 100%;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
}
</style>
