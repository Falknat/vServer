import { defineStore } from 'pinia'

export const useProxiesStore = defineStore('proxies', {
  state: () => ({
    list: [],
    loaded: false,
  }),

  actions: {
    async load() {
      try {
        const data = await api.getProxyList()
        if (data) {
          this.list = data
          this.loaded = true
        }
      } catch (e) {
        console.error('Failed to load proxies:', e)
      }
    },

    async create(proxyData) {
      try {
        const config = await api.getConfig()
        config.Proxy_Service.push({
          Name: proxyData.name || '',
          Enable: proxyData.enabled,
          ExternalDomain: proxyData.domain,
          LocalAddress: proxyData.localAddr,
          LocalPort: proxyData.localPort,
          ServiceHTTPSuse: proxyData.serviceHttps,
          AutoHTTPS: proxyData.autoHttps,
          Compression: proxyData.compression !== undefined ? proxyData.compression : true,
          AutoCreateSSL: proxyData.autoSSL,
        })
        const result = await api.saveConfig(JSON.stringify(config))
        if (isSuccess(result)) {
          await api.reloadConfig()
          await this.load()
        }
        return result
      } catch (e) {
        console.error('Failed to create proxy:', e)
        return `Error: ${e.message}`
      }
    },

    async remove(domain) {
      try {
        const config = await api.getConfig()
        config.Proxy_Service = config.Proxy_Service.filter(
          p => p.ExternalDomain !== domain
        )
        const result = await api.saveConfig(JSON.stringify(config))
        if (isSuccess(result)) {
          await api.reloadConfig()
          await this.load()
        }
        return result
      } catch (e) {
        console.error('Failed to remove proxy:', e)
        return `Error: ${e.message}`
      }
    },
  },
})
