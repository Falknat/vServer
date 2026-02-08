import { defineStore } from 'pinia'

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

    async create(proxyData) {
      const config = await api.getConfig()
      config.Proxy_Service.push({
        Name: proxyData.name || '',
        Enable: proxyData.enabled,
        ExternalDomain: proxyData.domain,
        LocalAddress: proxyData.localAddr,
        LocalPort: proxyData.localPort,
        ServiceHTTPSuse: proxyData.serviceHttps,
        AutoHTTPS: proxyData.autoHttps,
        AutoCreateSSL: proxyData.autoSSL,
      })
      const result = await api.saveConfig(JSON.stringify(config))
      if (result && !String(result).startsWith('Error')) {
        await api.reloadConfig()
        await this.load()
      }
      return result
    },

    async remove(domain) {
      const config = await api.getConfig()
      config.Proxy_Service = config.Proxy_Service.filter(
        p => p.ExternalDomain !== domain
      )
      const result = await api.saveConfig(JSON.stringify(config))
      if (result && !String(result).startsWith('Error')) {
        await api.reloadConfig()
        await this.load()
      }
      return result
    },
  },
})
