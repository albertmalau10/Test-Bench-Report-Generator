import { defineStore } from 'pinia'
import api from '../services/api' 

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: sessionStorage.getItem('token') || null,
    user: JSON.parse(sessionStorage.getItem('user')) || null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin'
  },
  actions: {
    async login(username, password) {
      const response = await api.post('/login', { username, password })
      const { token, user } = response.data

      this.token = token
      this.user = user

      sessionStorage.setItem('token', token)
      sessionStorage.setItem('user', JSON.stringify(user))
    },
    logout() {
      this.token = null
      this.user = null

      sessionStorage.removeItem('token')
      sessionStorage.removeItem('user')
    }
  }
})