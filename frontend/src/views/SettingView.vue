<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import api from "../services/api";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";
import Dropdown from "primevue/dropdown";
import Toast from "primevue/toast";
import { useToast } from "primevue/usetoast";

const toast = useToast();
const loading = ref(false);
let pollInterval = null;

const settings = ref({
  opc_ua_address: "",
  opc_node_pressure: "",
  opc_node_command: "",
  opc_node_feedback: "",
  opc_node_flow: "",
  opc_node_output: "",
  opc_ua_username: "",
  opc_ua_password: "",
  opc_ua_security_policy: "None",
  opc_ua_security_mode: "None",
});


const opcStatus = ref({
  server: false,
  command: { connected: false, value: "-" },
  feedback: { connected: false, value: "-" },
  pressure: { connected: false, value: "-" },
  flow: { connected: false, value: "-" },
  output: { connected: false, value: "-" },
});

const securityPolicies = ref([
  { label: "None", value: "None" },
  { label: "Basic256Sha256", value: "Basic256Sha256" },
  { label: "Aes128_Sha256_RsaOaep", value: "Aes128_Sha256_RsaOaep" },
  { label: "Aes256_Sha256_RsaPss", value: "Aes256_Sha256_RsaPss" },
]);

const securityModes = ref([
  { label: "None", value: "None" },
  { label: "Sign", value: "Sign" },
  { label: "Sign & Encrypt", value: "SignAndEncrypt" },
]);

async function fetchStatus() {
  try {
    const { data } = await api.get("/opcua/status");
    opcStatus.value = data;
  } catch (error) {
    opcStatus.value = {
      server: false,
      command: { connected: false, value: "-" },
      feedback: { connected: false, value: "-" },
      pressure: { connected: false, value: "-" },
      flow: { connected: false, value: "-" },
      output: { connected: false, value: "-" },
    };
  }
}

async function loadSettings() {
  try {
    const { data } = await api.get("/settings");
    settings.value.opc_ua_address = data.opc_ua_address || "";
    settings.value.opc_node_pressure = data.opc_node_pressure || "";
    settings.value.opc_node_command = data.opc_node_command || "";
    settings.value.opc_node_feedback = data.opc_node_feedback || "";
    settings.value.opc_node_flow = data.opc_node_flow || "";
    settings.value.opc_node_output = data.opc_node_output || "";
    settings.value.opc_ua_username = data.opc_ua_username || "";
    settings.value.opc_ua_security_policy =
      data.opc_ua_security_policy || "None";
    settings.value.opc_ua_security_mode = data.opc_ua_security_mode || "None";
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: "Failed to load settings",
      life: 3000,
    });
  }
}

async function saveSettings() {
  loading.value = true;
  try {
    await api.put("/settings", settings.value);
    toast.add({
      severity: "success",
      summary: "Success",
      detail: "Settings saved",
      life: 3000,
    });
    settings.value.opc_ua_password = "";
    await fetchStatus();
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: "Failed to save settings",
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadSettings();
  fetchStatus();
  pollInterval = setInterval(fetchStatus, 500);
});

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval);
});
</script>

<template>
  <div class="settings-page">
    <Toast />
    <h2>System Settings</h2>

    <div class="settings-card">
      <div class="section-title">
        <span>OPC UA Connection (ctrlX CORE)</span>
        <div class="server-status">
          <span class="status-text">{{
            opcStatus.server ? "Connected" : "Disconnected"
          }}</span>
          <span
            class="status-dot"
            :class="opcStatus.server ? 'connected' : 'disconnected'"
            title="Server Status"></span>
        </div>
      </div>

      <div class="field">
        <label>OPC-UA Address</label>
        <InputText
          v-model="settings.opc_ua_address"
          placeholder="e.g., opc.tcp://192.168.1.1:4840" />
      </div>

      <div class="node-grid">
        <div class="field">
          <label>Username</label>
          <InputText
            v-model="settings.opc_ua_username"
            placeholder="Leave empty for Anonymous" />
        </div>
        <div class="field">
          <label>Password (Update Only)</label>
          <Password
            v-model="settings.opc_ua_password"
            :feedback="false"
            toggleMask />
        </div>
        <div class="field">
          <label>Security Policy</label>
          <Dropdown
            v-model="settings.opc_ua_security_policy"
            :options="securityPolicies"
            optionLabel="label"
            optionValue="value" />
        </div>
        <div class="field">
          <label>Security Mode</label>
          <Dropdown
            v-model="settings.opc_ua_security_mode"
            :options="securityModes"
            optionLabel="label"
            optionValue="value" />
        </div>
      </div>

      <div class="section-title mt-4">OPC UA Node Identifiers</div>
      <div class="node-grid">
        <div class="field">
          <div class="label-row">
            <label>Output Trigger Node ID</label>
            <span
              class="status-dot"
              :class="opcStatus.output.connected ? 'connected' : 'disconnected'"
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.output.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_output"
              placeholder="ns=2;s=plc/app/Application/sym/PLC_PRG/output" />
          </div>
        </div>
        <div class="field">
          <div class="label-row">
            <label>Command Node ID</label>
            <span
              class="status-dot"
              :class="
                opcStatus.command.connected ? 'connected' : 'disconnected'
              "
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.command.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_command"
              placeholder="ns=2;s=testbench/sensors/command" />
          </div>
        </div>
        <div class="field">
          <div class="label-row">
            <label>Feedback Node ID</label>
            <span
              class="status-dot"
              :class="
                opcStatus.feedback.connected ? 'connected' : 'disconnected'
              "
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.feedback.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_feedback"
              placeholder="ns=2;s=testbench/sensors/feedback" />
          </div>
        </div>
        <div class="field">
          <div class="label-row">
            <label>Pressure Node ID</label>
            <span
              class="status-dot"
              :class="
                opcStatus.pressure.connected ? 'connected' : 'disconnected'
              "
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.pressure.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_pressure"
              placeholder="ns=2;s=testbench/sensors/pressure" />
          </div>
        </div>
        <div class="field">
          <div class="label-row">
            <label>Flow Node ID</label>
            <span
              class="status-dot"
              :class="opcStatus.flow.connected ? 'connected' : 'disconnected'"
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.flow.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_flow"
              placeholder="ns=2;s=testbench/sensors/flow" />
          </div>
        </div>
        <div class="field">
          <div class="label-row">
            <label>HPU Status Node ID</label>
            <span
              class="status-dot"
              :class="opcStatus.output.connected ? 'connected' : 'disconnected'"
              title="Node Status"></span>
          </div>
          <div class="input-with-value">
            <InputText
              :value="opcStatus.output.value"
              readonly
              class="live-value-box"
              tabindex="-1" />
            <InputText
              v-model="settings.opc_node_output"
              placeholder="ns=2;s=plc/app/Application/sym/PLC_PRG/output" />
          </div>
        </div>
      </div>

      <div class="actions">
        <Button
          label="Save Settings"
          icon="pi pi-save"
          :loading="loading"
          @click="saveSettings" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  padding: 1rem;
  overflow-y: auto;
  height: 100%;
}

h2 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-size: 1.5rem;
}

.settings-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  max-width: 1200px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: var(--primary-color);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.25rem;
  margin-bottom: -0.25rem;
}

.server-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8rem;
  color: var(--text-muted);
  font-weight: 500;
}

.mt-4 {
  margin-top: 0.5rem;
}

.node-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem 1.5rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.15rem; /* Tighter label spacing */
}

.field label {
  font-size: 0.78rem;
  color: var(--text-muted);
}

.label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.input-with-value {
  display: flex;
  gap: 0.5rem;
}

.input-with-value > :nth-child(2) {
  flex: 1;
}

/* Force PrimeVue inputs to be shorter */
:deep(.p-inputtext),
:deep(.p-dropdown),
:deep(.p-password input) {
  padding: 0.35rem 0.5rem;
  font-size: 0.85rem;
}

/* Restyled Readonly Value Box */
.live-value-box {
  width: 5rem;
  text-align: center;
  font-weight: 600;
  color: var(--primary-color);
  background-color: rgba(0, 0, 0, 0.03) !important;
  border: 1px solid var(--border-color);
}

.p-dark .live-value-box {
  background-color: rgba(255, 255, 255, 0.05) !important;
}

/* Status Indicator Dots */
.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: #ef4444;
  transition:
    background-color 0.3s ease,
    box-shadow 0.3s ease;
  flex-shrink: 0;
}
.status-dot.connected {
  background-color: #10b981;
  box-shadow: 0 0 6px rgba(16, 185, 129, 0.5);
}
.status-dot.disconnected {
  box-shadow: 0 0 6px rgba(239, 68, 68, 0.5);
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 0.25rem;
}

:deep(.p-button) {
  padding: 0.4rem 1rem;
}

@media (max-width: 600px) {
  .node-grid {
    grid-template-columns: 1fr;
  }
}
</style>
