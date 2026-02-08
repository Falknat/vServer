import { defineStore } from 'pinia'

const SERVICE_METHODS = {
  HTTP:  { start: () => api.startHTTPService(),  stop: () => api.stopHTTPService() },
  HTTPS: { start: () => api.startHTTPSService(), stop: () => api.stopHTTPSService() },
  MySQL: { start: () => api.startMySQLService(),  stop: () => api.stopMySQLService() },
  PHP:   { start: () => api.startPHPService(),   stop: () => api.stopPHPService() },
}

export const useServicesStore = defineStore('services', {
  state: () => ({
    list: [],
    loaded: false,
    isOperating: false,
  }),

  actions: {
    async load() {
      try {
        const data = await api.getAllServicesStatus()
        if (data) {
          this.list = Array.isArray(data) ? data : Object.values(data)
          this.loaded = true
        }
      } catch (e) {
        console.error('Failed to load services:', e)
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

    async toggleService(name, action) {
      try {
        const method = SERVICE_METHODS[name]?.[action]
        if (method) await method()
        await this.load()
      } catch (e) {
        console.error(`Failed to ${action} service ${name}:`, e)
      }
    },

    async startService(name) {
      return this.toggleService(name, 'start')
    },

    async stopService(name) {
      return this.toggleService(name, 'stop')
    },

    async enableProxy() {
      try {
        return await api.enableProxyService()
      } catch (e) {
        console.error('Failed to enable proxy:', e)
      }
    },

    async disableProxy() {
      try {
        return await api.disableProxyService()
      } catch (e) {
        console.error('Failed to disable proxy:', e)
      }
    },

    async enableACME() {
      try {
        return await api.enableACMEService()
      } catch (e) {
        console.error('Failed to enable ACME:', e)
      }
    },

    async disableACME() {
      try {
        return await api.disableACMEService()
      } catch (e) {
        console.error('Failed to disable ACME:', e)
      }
    },

    async restartAll() {
      try {
        return await api.restartAllServices()
      } catch (e) {
        console.error('Failed to restart services:', e)
      }
    },
  },
})
