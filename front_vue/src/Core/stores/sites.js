import { defineStore } from 'pinia'
import { api } from '@core/api/index.js'

export const useSitesStore = defineStore('sites', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async load() {
      const data = await api.getSitesList()
      if (data) {
        this.list = data
        this.loaded = true
      }
    },

    async create(siteData) {
      const result = await api.createNewSite(JSON.stringify(siteData))
      if (result && !String(result).startsWith('Error')) {
        await this.load()
      }
      return result
    },

    async remove(host) {
      const result = await api.deleteSite(host)
      if (result && !String(result).startsWith('Error')) {
        await this.load()
      }
      return result
    },

    async openFolder(host) {
      await api.openSiteFolder(host)
    },

    async uploadCert(host, certType, certDataBase64) {
      return await api.uploadCertificate(host, certType, certDataBase64)
    },
  },
})
