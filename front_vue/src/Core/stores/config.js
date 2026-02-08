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
      try {
        const config = await api.getConfig()
        if (config) {
          this.data = config
          this.loaded = true
        }
      } catch (e) {
        console.error('Failed to load config:', e)
      }
    },

    async save(configData) {
      try {
        const result = await api.saveConfig(JSON.stringify(configData, null, 4))
        if (isSuccess(result)) {
          this.data = configData
        }
        return result
      } catch (e) {
        console.error('Failed to save config:', e)
        return `Error: ${e.message}`
      }
    },
  },
})
