<template>
  <div class="dashboard-view">
    <div class="welcome-container">
      <header class="welcome-header">
        <h1>Test Bench Report Generator</h1>
        <p>Welcome back. Select a module to begin testing or managing data.</p>
      </header>
      
      <div class="action-cards">
        <div class="action-card" @click="goToRecords">
          <div class="icon-wrapper">
            <i class="pi pi-chart-line"></i>
          </div>
          <h3>Test Workspace</h3>
          <p>Monitor live OPC UA telemetry, actuate valves, and record test data performance.</p>
          <Button label="Open Workspace" icon="pi pi-arrow-right" severity="secondary" outlined class="action-btn" />
        </div>

        <div class="action-card" @click="goToValveList">
          <div class="icon-wrapper">
            <i class="pi pi-server"></i>
          </div>
          <h3>Valve Database</h3>
          <p>Manage the master directory of registered valves, specifications, and PDF datasheets.</p>
          <Button label="View Database" icon="pi pi-arrow-right" severity="secondary" outlined class="action-btn" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Button from 'primevue/button'

const router = useRouter()

function goToRecords() {
  router.push('/records')
}

function goToValveList() {
  router.push('/valves')
}
</script>

<style scoped>
.dashboard-view {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  overflow-y: auto;
}

.welcome-container {
  max-width: 900px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 3.5rem;
  text-align: center;
}

.welcome-header h1 {
  font-size: 2.5rem;
  margin: 0 0 0.75rem 0;
  color: var(--text-color);
  letter-spacing: -0.03em;
}

.welcome-header p {
  font-size: 1.1rem;
  color: var(--text-muted);
  margin: 0;
}

.action-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.action-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 3rem 2.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
  border-color: var(--primary-color);
}

.action-card:hover .icon-wrapper {
  background: var(--primary-color);
  color: #fff;
}

.icon-wrapper {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: rgba(0, 43, 73, 0.05);
  color: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.5rem;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.p-dark .icon-wrapper {
  background: rgba(255, 255, 255, 0.05);
}

.icon-wrapper i {
  font-size: 2rem;
}

.action-card h3 {
  margin: 0 0 1rem 0;
  font-size: 1.4rem;
  color: var(--text-color);
}

.action-card p {
  margin: 0 0 2rem 0;
  color: var(--text-muted);
  line-height: 1.5;
  font-size: 0.95rem;
}

.action-btn {
  margin-top: auto;
  pointer-events: none; /* Let the parent card handle the click entirely */
}

@media (max-width: 768px) {
  .action-cards {
    grid-template-columns: 1fr;
  }
}
</style>