import { defineStore } from 'pinia'
import { api } from '@core/api/index.js'

export const useProxiesStore = defineStore('proxies', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async load() {
      const data = await api.getProxyList()
      if (data) {
        this.list = data
        this.loaded = true
      }
    },
  },
})
