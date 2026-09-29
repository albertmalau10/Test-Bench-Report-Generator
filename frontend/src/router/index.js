import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'

import DashboardView from '../views/DashboardView.vue'
import ValveListView from '../views/ValveListView.vue'
import ValveDetailsView from '../views/ValveDetailsView.vue'
import RecordsView from '../views/RecordsView.vue'
import SettingView from '../views/SettingView.vue'
import LoginView from '../views/LoginView.vue'
import { useAuthStore } from '../stores/auth'



const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
  path: '/setting',
  name: 'setting',
  component: () => import('../views/SettingView.vue'),
  meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'dashboard', component: DashboardView },
      { path: 'valves', name: 'valve-list', component: ValveListView },
      { path: 'valves/details', name: 'valve-details-search', component: ValveDetailsView },
      { path: 'valves/:id', name: 'valve-details', component: ValveDetailsView },
      { path: 'records', name: 'records', component: RecordsView },
      { path: 'settings', name: 'settings', component: SettingView, meta: { requiresAdmin: true } },
    ]
  }
]

const router = createRouter({
  // import.meta.env.BASE_URL follows Vite's `base` config (see
  // vite.config.js) so client-side navigation still works once the
  // app is served under /valve-database-app/ behind ctrlX CORE's
  // reverse proxy, not just at the site root.
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'login' })
  } else if (to.name === 'login' && authStore.isLoggedIn) {
    next({ name: 'dashboard' })
  // FIX: Enforce admin role check
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' }) 
  } else {
    next()
  }
})

export default router