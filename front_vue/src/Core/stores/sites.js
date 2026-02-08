import { defineStore } from 'pinia'

export const useSitesStore = defineStore('sites', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async load() {
      try {
        const data = await api.getSitesList()
        if (data) {
          this.list = data
          this.loaded = true
        }
      } catch (e) {
        console.error('Failed to load sites:', e)
      }
    },

    async create(siteData) {
      try {
        const result = await api.createNewSite(JSON.stringify(siteData))
        if (isSuccess(result)) {
          await this.load()
        }
        return result
      } catch (e) {
        console.error('Failed to create site:', e)
        return `Error: ${e.message}`
      }
    },

    async remove(host) {
      try {
        const result = await api.deleteSite(host)
        if (isSuccess(result)) {
          await this.load()
        }
        return result
      } catch (e) {
        console.error('Failed to delete site:', e)
        return `Error: ${e.message}`
      }
    },

    async openFolder(host) {
      try {
        await api.openSiteFolder(host)
      } catch (e) {
        console.error('Failed to open folder:', e)
      }
    },

    async uploadCert(host, certType, certDataBase64) {
      try {
        return await api.uploadCertificate(host, certType, certDataBase64)
      } catch (e) {
        console.error('Failed to upload cert:', e)
        return `Error: ${e.message}`
      }
    },
  },
})
