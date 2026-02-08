import { defineStore } from 'pinia'

export const useServicesStore = defineStore('services', {
  state: () => ({
    list: [],
    loaded: false,
    isOperating: false,
  }),

  actions: {
    async load() {
      const data = await api.getAllServicesStatus()
      if (data) {
        this.list = Array.isArray(data) ? data : Object.values(data)
        this.loaded = true
      }
    },

    setPending(text) {
      this.isOperating = true
      this.list = this.list.map(s => ({ ...s, pending: text }))
    },

    async clearPending() {
      await this.load()
      this.isOperating = false
    },

    async startService(name) {
      const methods = {
        HTTP: () => api.startHTTPService(),
        HTTPS: () => api.startHTTPSService(),
        MySQL: () => api.startMySQLService(),
        PHP: () => api.startPHPService(),
      }
      if (methods[name]) await methods[name]()
      await this.load()
    },

    async stopService(name) {
      const methods = {
        HTTP: () => api.stopHTTPService(),
        HTTPS: () => api.stopHTTPSService(),
        MySQL: () => api.stopMySQLService(),
        PHP: () => api.stopPHPService(),
      }
      if (methods[name]) await methods[name]()
      await this.load()
    },
  },
})
