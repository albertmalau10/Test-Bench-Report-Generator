import { defineStore } from 'pinia'

export const useValveStore = defineStore('valve', {
  state: () => ({
    selectedValve: null
  }),
  actions: {
    selectValve(valve) {
      this.selectedValve = valve
    }
  }
})