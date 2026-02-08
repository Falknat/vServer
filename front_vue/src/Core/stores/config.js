import { defineStore } from 'pinia'

export const useConfigStore = defineStore('config', {
  state: () => ({
    data: null,
    loaded: false,
  }),

  getters: {
    softSettings: (state) => state.data?.Soft_Settings || {},
  },

  actions: {
    async load() {
      const config = await api.getConfig()
      if (config) {
        this.data = config
        this.loaded = true
      }
    },

    async save(configData) {
      const result = await api.saveConfig(JSON.stringify(configData, null, 4))
      if (result && !String(result).startsWith('Error')) {
        this.data = configData
      }
      return result
    },

    async enableProxy() {
      return await api.enableProxyService()
    },

    async disableProxy() {
      return await api.disableProxyService()
    },

    async enableACME() {
      return await api.enableACMEService()
    },

    async disableACME() {
      return await api.disableACMEService()
    },

    async restartAll() {
      return await api.restartAllServices()
    },
  },
})
