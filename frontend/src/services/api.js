import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// The backend for this build lives at <BASE_URL>api — e.g.
// "/valve-database-app/api" once behind ctrlX CORE's reverse proxy
// (see vite.config.js's `base` and backend/main.go's /api route
// group), or "/api" for a plain local `npm run dev` at "/". Set
// VITE_API_URL to point at a backend on a different origin instead
// (e.g. http://localhost:8080/api for local dev against `go run .`
// directly, without the Vite dev-server proxy).
export const API_ORIGIN = import.meta.env.VITE_API_URL || `${import.meta.env.BASE_URL}api`

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const isLoginRequest = error.config?.url?.includes('/login')

    if (error.response?.status === 401 && !isLoginRequest) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export function uploadValveImage(id, file) {
  const formData = new FormData()
  formData.append('image', file)

  return api.post(`/valves/${id}/image`, formData, {
    headers: { 'Content-Type': undefined }
  })
}

export function uploadValveDatasheet(id, file) {
  const formData = new FormData()
  formData.append('datasheet', file)

  return api.post(`/valves/${id}/datasheet`, formData, {
    headers: { 'Content-Type': undefined }
  })
}

export function startOutput() {
  return api.post('/ctrlx/start') // Route kept same for minimal disruption
}

export function stopOutput() {
  return api.post('/ctrlx/stop')
}

export function fetchOpcData() {
  return api.get('/opcua/data')
}

export default api
