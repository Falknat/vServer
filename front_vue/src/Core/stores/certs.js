import { defineStore } from 'pinia'
import { api } from '@core/api/index.js'

export const useCertsStore = defineStore('certs', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async loadAll() {
      const data = await api.getAllCertsInfo()
      if (data) {
        this.list = data
        this.loaded = true
      }
    },

    async getInfo(domain) {
      return await api.getCertInfo(domain)
    },

    async issue(domain) {
      const result = await api.obtainSSLCertificate(domain)
      await this.loadAll()
      return result
    },

    async renew(domain) {
      const result = await api.obtainSSLCertificate(domain)
      await this.loadAll()
      return result
    },

    async remove(domain) {
      const result = await api.deleteCertificate(domain)
      await this.loadAll()
      return result
    },

    async reload() {
      return await api.reloadSSLCertificates()
    },
  },
})
