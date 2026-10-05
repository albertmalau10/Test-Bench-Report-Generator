<template>
  <main class="login-page">
    <section class="login-shell">
      <!-- Left Identity Panel -->
      <aside class="system-panel">
        <div class="system-brand">
          <div class="brand-badge">
            <img :src="rexrothLogo" alt="Bosch Rexroth" class="brand-logo" />
          </div>
          <div class="brand-copy">
            <span class="brand-label">PLATFORM</span>
            <strong>ValveDAX</strong>
          </div>
        </div>

        <div class="system-content">
          <p class="section-index">01 / LOGIN PAGE</p>
          <h1>Valve Data Acquisition</h1>
          <p class="system-description">
            Dedicated hydraulic valve data acquisition stand powered by Bosch Rexroth ctrlX CORE.
          </p>

          <dl class="system-information">
            <div class="information-row">
              <dt>System</dt>
              <dd>ctrlX CORE Bench</dd>
            </div>
            <div class="information-row">
              <dt>Protocol</dt>
              <dd>OPC UA (TCP:4840)</dd>
            </div>
            <div class="information-row">
              <dt>Access</dt>
              <dd>
                <span class="status-indicator"></span>
                Local Air-Gapped
              </dd>
            </div>
          </dl>
        </div>

        <div class="system-footer">
          <span>Bosch Rexroth 2026</span>
          <span>v1.0.0</span>
        </div>
      </aside>

      <!-- Right Form Panel -->
      <section class="form-panel">
        <div class="login-form-container">
          <header class="form-header">
            <span class="form-code">AUTH / 01</span>
            <h2>Sign In</h2>
            <p>Enter assigned workshop operator credentials.</p>
          </header>

          <form class="login-form" novalidate @submit.prevent="handleLogin">
            <div class="field">
              <label for="username">Username</label>
              <div class="control-wrapper">
                <i class="pi pi-user control-icon"></i>
                <InputText
                  id="username"
                  v-model.trim="username"
                  type="text"
                  autocomplete="username"
                  placeholder="Enter username (e.g. operator)"
                  :disabled="isLoading"
                  :invalid="Boolean(fieldErrors.username)"
                  class="text-control"
                  autofocus
                  @input="clearError('username')"
                />
              </div>
              <small v-if="fieldErrors.username" class="field-error">
                {{ fieldErrors.username }}
              </small>
            </div>

            <div class="field">
              <label for="password">Password</label>
              <div class="control-wrapper">
                <i class="pi pi-lock control-icon"></i>
                <Password
                  v-model="password"
                  inputId="password"
                  autocomplete="current-password"
                  placeholder="Enter password"
                  :feedback="false"
                  :disabled="isLoading"
                  :invalid="Boolean(fieldErrors.password)"
                  toggleMask
                  class="password-control"
                  inputClass="password-input"
                  @input="clearError('password')"
                />
              </div>
              <small v-if="fieldErrors.password" class="field-error">
                {{ fieldErrors.password }}
              </small>
            </div>

            <div v-if="errorMessage" class="authentication-error" role="alert">
              <i class="pi pi-exclamation-triangle"></i>
              <div>
                <strong>Authentication failed</strong>
                <span>{{ errorMessage }}</span>
              </div>
            </div>

            <Button
              type="submit"
              label="Login"
              icon="pi pi-arrow-right"
              iconPos="right" 
              :loading="isLoading"
              :disabled="isLoading"
              class="login-button"
            />
          </form>

          <!-- Quick Operator Credentials Guide -->
          <footer class="form-footer">
            <div class="credential-pills">
              <span class="cred-chip"><strong>Role:</strong> operator</span>
              <span class="cred-chip"><strong>Admin:</strong> admin</span>
            </div>
          </footer>
        </div>
      </section>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

import rexrothLogo from '../assets/Bosch_Rexroth-Logo.wine.svg'

import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";

const username = ref("");
const password = ref("");
const errorMessage = ref("");
const isLoading = ref(false);

const fieldErrors = reactive({
  username: "",
  password: "",
});

const authStore = useAuthStore();
const router = useRouter();

function validateForm() {
  fieldErrors.username = "";
  fieldErrors.password = "";
  errorMessage.value = "";

  if (!username.value.trim()) {
    fieldErrors.username = "Username is required.";
  }

  if (!password.value) {
    fieldErrors.password = "Password is required.";
  }

  return !fieldErrors.username && !fieldErrors.password;
}

function clearError(fieldName) {
  fieldErrors[fieldName] = "";
  errorMessage.value = "";
}

async function handleLogin() {
  if (!validateForm()) {
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    await authStore.login(username.value.trim(), password.value);

    await router.push("/");
  } catch (err) {
    if (!err.response) {
      errorMessage.value = "Unable to connect to the backend server.";
    } else if (err.response.status === 401) {
      errorMessage.value = "Incorrect username or password.";
    } else {
      errorMessage.value =
        err.response?.data?.error || "An unexpected error occurred.";
    }
  } finally {
    isLoading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  width: 100%;
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 2rem;
  color: var(--text-color);
  background-color: var(--bg-color);
  background-image:
    linear-gradient(rgba(0, 43, 73, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 43, 73, 0.05) 1px, transparent 1px);
  background-size: 32px 32px;
}

.p-dark .login-page {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
}

.login-shell {
  width: min(100%, 960px);
  min-height: 610px;
  display: grid;
  grid-template-columns: 42% 58%;
  overflow: hidden;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.3);
}

/* Left system panel */
.system-panel {
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 2.25rem;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
}

.system-brand {
  display: flex;
  align-items: center;
  gap: 0.9rem;
}

.logo-frame {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  background: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}

.brand-logo {
  height: 22px;
  width: auto;
  object-fit: contain;
  display: block;
}

.brand-badge {
  background: #ffffff;
  padding: 0.35rem 0.65rem;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.brand-copy {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.brand-label {
  color: var(--rexroth-casper);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.65rem;
  letter-spacing: 0.14em;
}

.brand-copy strong {
  color: #ffffff;
  font-size: 1.05rem;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.system-brand .brand-label {
  color: var(--text-muted);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.62rem;
  letter-spacing: 0.16em;
}

.system-brand .brand-copy strong {
  color: #ffffff;
  font-size: 0.88rem;
  font-weight: 600;
  line-height: 1.35;
}

.mobile-brand .brand-label {
  color: var(--text-muted);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.62rem;
  letter-spacing: 0.16em;
}

.mobile-brand .brand-copy strong {
  color: var(--text-color);
  font-size: 0.88rem;
  font-weight: 600;
  line-height: 1.35;
}

.system-content {
  margin: auto 0;
  padding: 3rem 0;
}

.section-index {
  margin: 0 0 1rem;
  color: var(--info-color);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: 0.13em;
}

.system-content h1 {
  max-width: 22rem;
  margin: 0;
  color: #ffffff;
  font-size: clamp(2rem, 3.2vw, 2.8rem);
  font-weight: 600;
  line-height: 1.12;
  letter-spacing: -0.035em;
}

.system-description {
  max-width: 21rem;
  margin: 1.25rem 0 0;
  color: var(--rexroth-casper);
  font-size: 0.86rem;
  line-height: 1.75;
}

.system-information {
  margin: 2.6rem 0 0;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.information-row {
  display: grid;
  grid-template-columns: 6rem 1fr;
  min-height: 2.7rem;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.information-row dt {
  color: var(--rexroth-casper);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.information-row dd {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.48rem;
  color: #ffffff;
  font-size: 0.78rem;
}

.status-indicator {
  width: 7px;
  height: 7px;
  background: var(--info-color);
  border-radius: 50%;
}

.system-footer {
  display: flex;
  justify-content: space-between;
  color: var(--rexroth-casper);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.65rem;
  letter-spacing: 0.04em;
}

/* Right form panel */
.form-panel {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3.5rem;
  background: var(--card-bg);
}

.mobile-brand {
  display: none;
}

.login-form-container {
  width: 100%;
  max-width: 340px;
}

.form-header {
  margin-bottom: 2.25rem;
}

.form-code {
  margin: 0 0 0.85rem;
  color: var(--info-color);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: 0.13em;
}

.form-header h2 {
  margin: 0 0 0.5rem;
  color: var(--text-color);
  font-size: 1.75rem;
  font-weight: 600;
  letter-spacing: -0.025em;
}

.form-header p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.84rem;
  line-height: 1.6;
}

.form-footer {
  margin-top: 1.75rem;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding-top: 1rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.35rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.48rem;
}

.field label {
  color: var(--text-color);
  font-size: 0.77rem;
  font-weight: 600;
}

.credential-pills {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.cred-chip {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 0.2rem 0.55rem;
  border-radius: 4px;
  font-family: Consolas, Monaco, monospace;
  font-size: 0.7rem;
  color: var(--text-muted);
}

.cred-chip strong {
  color: var(--primary-color);
}

.control-wrapper {
  position: relative;
  width: 100%;
}

.control-icon {
  position: absolute;
  z-index: 2;
  top: 50%;
  left: 0.95rem;
  color: var(--text-muted);
  font-size: 0.86rem;
  pointer-events: none;
  transform: translateY(-50%);
}

.control-wrapper:focus-within .control-icon {
  color: var(--primary-color);
}

.text-control {
  width: 100%;
}

.control-wrapper :deep(.p-inputtext),
.control-wrapper :deep(.password-input) {
  width: 100%;
  height: 46px;
  padding-left: 2.65rem;
  color: var(--text-color);
  background: var(--bg-color);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 0.85rem;
  box-shadow: none;
  transition:
    border-color 0.15s ease,
    background-color 0.15s ease;
}

.control-wrapper :deep(.p-inputtext::placeholder),
.control-wrapper :deep(.password-input::placeholder) {
  color: var(--text-muted);
}

.control-wrapper :deep(.p-inputtext:hover),
.control-wrapper :deep(.password-input:hover) {
  border-color: var(--primary-color);
}

.control-wrapper :deep(.p-inputtext:focus),
.control-wrapper :deep(.password-input:focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(0, 204, 255, 0.2);
}

.control-wrapper :deep(.p-invalid) {
  border-color: var(--danger-color);
}

.password-control {
  width: 100%;
}

.password-control :deep(.p-password) {
  width: 100%;
}

.password-control :deep(.password-input) {
  padding-right: 2.8rem;
}

.password-control :deep(.p-password-toggle-mask-icon) {
  color: var(--text-muted);
}

.field-error {
  color: var(--danger-color);
  font-size: 0.7rem;
}

.authentication-error {
  display: flex;
  align-items: flex-start;
  gap: 0.7rem;
  padding: 0.85rem;
  color: #ff8080;
  background: rgba(223, 0, 36, 0.1);
  border-left: 3px solid var(--danger-color);
  border-radius: 2px;
}

.authentication-error > i {
  margin-top: 0.12rem;
  color: var(--danger-color);
  font-size: 0.85rem;
}

.authentication-error > div {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.authentication-error strong {
  color: var(--danger-color);
  font-size: 0.75rem;
  font-weight: 600;
}

.authentication-error span {
  font-size: 0.7rem;
  line-height: 1.45;
}

.login-button {
  width: 100%;
  height: 46px;
  margin-top: 0.2rem;
  display: flex;
  justify-content: space-between;
  color: #001524;
  background: var(--info-color);
  border: 1px solid var(--info-color);
  border-radius: 4px;
  font-size: 0.82rem;
  font-weight: 700;
  box-shadow: none;
  transition:
    background-color 0.15s ease,
    border-color 0.15s ease;
}

.login-button:not(:disabled):hover {
  background: #33d6ff;
  border-color: #33d6ff;
}

.form-footer {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  margin-top: 1.5rem;
  color: var(--text-muted);
  font-size: 0.68rem;
}

.form-footer i {
  color: var(--text-muted);
  font-size: 0.72rem;
}

@media (max-width: 760px) {
  .login-page {
    display: block;
    min-height: 100vh;
    padding: 0;
    overflow-y: auto;
  }

  .login-shell {
    display: block;
    min-height: 100vh;
    border: none;
    border-radius: 0;
    box-shadow: none;
  }

  .system-panel {
    display: none;
  }

  .form-panel {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 2rem;
  }

  .mobile-brand {
    width: 100%;
    max-width: 340px;
    display: flex;
    align-items: center;
    gap: 0.8rem;
    margin-bottom: 3.5rem;
  }

  .mobile-logo {
    width: 38px;
    height: 38px;
  }
}

@media (max-width: 420px) {
  .form-panel {
    justify-content: flex-start;
    padding: 1.5rem;
  }

  .mobile-brand {
    margin: 1rem 0 4rem;
  }

  .form-header h2 {
    font-size: 1.55rem;
  }
}
</style>
