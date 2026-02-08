import { defineStore } from 'pinia'

export const useCertsStore = defineStore('certs', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async loadAll() {
      try {
        const data = await api.getAllCertsInfo()
        if (data) {
          this.list = data
          this.loaded = true
        }
      } catch (e) {
        console.error('Failed to load certs:', e)
      }
    },

    async getInfo(domain) {
      try {
        return await api.getCertInfo(domain)
      } catch (e) {
        console.error('Failed to get cert info:', e)
        return null
      }
    },

    async _obtainAndReload(domain) {
      try {
        const result = await api.obtainSSLCertificate(domain)
        await this.loadAll()
        return result
      } catch (e) {
        console.error('Failed to obtain certificate:', e)
        return `Error: ${e.message}`
      }
    },

    async issue(domain) {
      return this._obtainAndReload(domain)
    },

    async renew(domain) {
      return this._obtainAndReload(domain)
    },

    async remove(domain) {
      try {
        const result = await api.deleteCertificate(domain)
        await this.loadAll()
        return result
      } catch (e) {
        console.error('Failed to delete certificate:', e)
        return `Error: ${e.message}`
      }
    },

    async reload() {
      try {
        return await api.reloadSSLCertificates()
      } catch (e) {
        console.error('Failed to reload certificates:', e)
      }
    },
  },
})
