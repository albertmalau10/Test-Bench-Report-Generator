<template>
  <main class="login-page">
    <section class="login-shell">
      <!-- System identity -->
      <aside class="system-panel">
        <div class="system-brand">
          <div class="logo-frame">
            <img
              :src="logo"
              alt="Test Bench Report Generator"
              class="brand-logo"
            />
          </div>

          <div class="brand-copy">
            <span class="brand-label">APPLICATION</span>
            <strong>Test Bench Report Generator</strong>
          </div>
        </div>

        <div class="system-content">
          <p class="section-index">01 / ACCESS</p>

          <h1>Engineering test data and report management.</h1>

          <p class="system-description">
            Authorized access for hydraulic valve testing, data recording,
            and report generation.
          </p>

          <dl class="system-information">
            <div class="information-row">
              <dt>System</dt>
              <dd>Test Bench Platform</dd>
            </div>

            <div class="information-row">
              <dt>Interface</dt>
              <dd>Web Application</dd>
            </div>

            <div class="information-row">
              <dt>Access</dt>
              <dd>
                <span class="status-indicator"></span>
                Restricted
              </dd>
            </div>
          </dl>
        </div>

        <div class="system-footer">
          <span>Internal Engineering Tool</span>
          <span>v1.0</span>
        </div>
      </aside>

      <!-- Authentication form -->
      <section class="form-panel">
        <div class="mobile-brand">
          <div class="logo-frame mobile-logo">
            <img
              :src="logo"
              alt="Test Bench Report Generator"
              class="brand-logo"
            />
          </div>

          <div class="brand-copy">
            <span class="brand-label">APPLICATION</span>
            <strong>Test Bench Report Generator</strong>
          </div>
        </div>

        <div class="login-form-container">
          <header class="form-header">
            <span class="form-code">AUTH / 01</span>
            <h2>Sign in</h2>
            <p>Enter your assigned account credentials.</p>
          </header>

          <form
            class="login-form"
            novalidate
            @submit.prevent="handleLogin"
          >
            <div class="field">
              <label for="username">Username</label>

              <div class="control-wrapper">
                <i class="pi pi-user control-icon"></i>

                <InputText
                  id="username"
                  v-model.trim="username"
                  type="text"
                  autocomplete="username"
                  placeholder="Enter username"
                  :disabled="isLoading"
                  :invalid="Boolean(fieldErrors.username)"
                  class="text-control"
                  autofocus
                  @input="clearError('username')"
                />
              </div>

              <small
                v-if="fieldErrors.username"
                class="field-error"
              >
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

              <small
                v-if="fieldErrors.password"
                class="field-error"
              >
                {{ fieldErrors.password }}
              </small>
            </div>

            <div
              v-if="errorMessage"
              class="authentication-error"
              role="alert"
            >
              <i class="pi pi-exclamation-triangle"></i>

              <div>
                <strong>Authentication failed</strong>
                <span>{{ errorMessage }}</span>
              </div>
            </div>

            <Button
              type="submit"
              label="Continue"
              icon="pi pi-arrow-right"
              iconPos="right"
              :loading="isLoading"
              :disabled="isLoading"
              class="login-button"
            />
          </form>

          <footer class="form-footer">
            <i class="pi pi-shield"></i>
            <span>Authorized personnel only</span>
          </footer>
        </div>
      </section>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

import logo from '../assets/logo.jpg'

import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

const fieldErrors = reactive({
  username: '',
  password: ''
})

const authStore = useAuthStore()
const router = useRouter()

function validateForm() {
  fieldErrors.username = ''
  fieldErrors.password = ''
  errorMessage.value = ''

  if (!username.value.trim()) {
    fieldErrors.username = 'Username is required.'
  }

  if (!password.value) {
    fieldErrors.password = 'Password is required.'
  }

  return !fieldErrors.username && !fieldErrors.password
}

function clearError(fieldName) {
  fieldErrors[fieldName] = ''
  errorMessage.value = ''
}

async function handleLogin() {
  if (!validateForm()) {
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.login(
      username.value.trim(),
      password.value
    )

    await router.push('/')
  } catch (err) {
    if (!err.response) {
      errorMessage.value =
        'Unable to connect to the backend server.'
    } else if (err.response.status === 401) {
      errorMessage.value =
        'Incorrect username or password.'
    } else {
      errorMessage.value =
        err.response?.data?.error ||
        'An unexpected error occurred.'
    }
  } finally {
    isLoading.value = false
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
  color: var(--text-primary);
  background-color: var(--bg-main);
  background-image:
    linear-gradient(
      rgba(0, 43, 73, 0.035) 1px,
      transparent 1px
    ),
    linear-gradient(
      90deg,
      rgba(0, 43, 73, 0.035) 1px,
      transparent 1px
    );
  background-size: 32px 32px;
}

.login-shell {
  width: min(100%, 960px);
  min-height: 610px;
  display: grid;
  grid-template-columns: 42% 58%;
  overflow: hidden;
  background: var(--bg-card);
  border: 1px solid var(--stopped);
  border-radius: 4px;
  box-shadow: 0 24px 60px rgba(0, 43, 73, 0.14);
}

/* Left system panel */

.system-panel {
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 2.25rem;
  background: var(--bosch-blue);
  border-right: 1px solid var(--bosch-blue-2);
}

.system-brand {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.logo-frame {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  background: #ffffff;
  border: 1px solid var(--bosch-blue-2);
  border-radius: 3px;
}

.brand-logo {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.brand-copy {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.22rem;
}

/* Brand text sits on two different backgrounds: the navy
   system-panel (desktop) and the white form-panel (mobile, <760px,
   see .mobile-brand below) — each needs its own contrast pair
   rather than one shared color. */
.system-brand .brand-label {
  color: var(--stopped);
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
  color: var(--text-secondary);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.62rem;
  letter-spacing: 0.16em;
}

.mobile-brand .brand-copy strong {
  color: var(--text-primary);
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
  color: var(--info);
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
  color: var(--stopped);
  font-size: 0.86rem;
  line-height: 1.75;
}

.system-information {
  margin: 2.6rem 0 0;
  border-top: 1px solid var(--bosch-blue-2);
}

.information-row {
  display: grid;
  grid-template-columns: 6rem 1fr;
  min-height: 2.7rem;
  align-items: center;
  border-bottom: 1px solid var(--bosch-blue-2);
}

.information-row dt {
  color: var(--stopped);
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
  background: var(--info);
  border-radius: 50%;
}

.system-footer {
  display: flex;
  justify-content: space-between;
  color: var(--stopped);
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
  background: #ffffff;
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
  color: var(--bosch-blue-2);
  font-family: Consolas, Monaco, monospace;
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: 0.13em;
}

.form-header h2 {
  margin: 0 0 0.5rem;
  color: var(--text-primary);
  font-size: 1.75rem;
  font-weight: 600;
  letter-spacing: -0.025em;
}

.form-header p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.84rem;
  line-height: 1.6;
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
  color: var(--text-primary);
  font-size: 0.77rem;
  font-weight: 600;
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
  color: var(--text-secondary);
  font-size: 0.86rem;
  pointer-events: none;
  transform: translateY(-50%);
}

.control-wrapper:focus-within .control-icon {
  color: var(--info);
}

.text-control {
  width: 100%;
}

.control-wrapper :deep(.p-inputtext),
.control-wrapper :deep(.password-input) {
  width: 100%;
  height: 46px;
  padding-left: 2.65rem;
  color: var(--text-primary);
  background: #ffffff;
  border: 1px solid var(--stopped);
  border-radius: 3px;
  font-size: 0.85rem;
  box-shadow: none;
  transition:
    border-color 0.15s ease,
    background-color 0.15s ease;
}

.control-wrapper :deep(.p-inputtext::placeholder),
.control-wrapper :deep(.password-input::placeholder) {
  color: var(--text-secondary);
}

.control-wrapper :deep(.p-inputtext:hover),
.control-wrapper :deep(.password-input:hover) {
  border-color: var(--bosch-blue-2);
}

.control-wrapper :deep(.p-inputtext:focus),
.control-wrapper :deep(.password-input:focus) {
  background: #ffffff;
  border-color: var(--info);
  box-shadow: 0 0 0 2px rgba(0, 200, 247, 0.16);
}

.control-wrapper :deep(.p-invalid) {
  border-color: var(--alarm-high);
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
  color: var(--text-secondary);
}

.field-error {
  color: var(--alarm-high);
  font-size: 0.7rem;
}

.authentication-error {
  display: flex;
  align-items: flex-start;
  gap: 0.7rem;
  padding: 0.85rem;
  color: #8a2020;
  background: #fdecec;
  border-left: 3px solid var(--alarm-high);
  border-radius: 2px;
}

.authentication-error > i {
  margin-top: 0.12rem;
  color: var(--alarm-high);
  font-size: 0.85rem;
}

.authentication-error > div {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.authentication-error strong {
  color: var(--alarm-high);
  font-size: 0.75rem;
  font-weight: 600;
}

.authentication-error span {
  color: #8a2020;
  font-size: 0.7rem;
  line-height: 1.45;
}

.login-button {
  width: 100%;
  height: 46px;
  margin-top: 0.2rem;
  display: flex;
  justify-content: space-between;
  color: var(--bosch-blue);
  background: var(--info);
  border: 1px solid var(--info);
  border-radius: 3px;
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

.login-button:not(:disabled):active {
  background: #00a8d1;
  border-color: #00a8d1;
}

.form-footer {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  margin-top: 1.5rem;
  color: var(--text-secondary);
  font-size: 0.68rem;
}

.form-footer i {
  color: var(--text-secondary);
  font-size: 0.72rem;
}

/* Responsive */

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
