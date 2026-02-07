import servicesData from './mock-data/services.json'
import sitesData from './mock-data/sites.json'
import proxiesData from './mock-data/proxies.json'
import configData from './mock-data/config.json'
import certsData from './mock-data/certs.json'
import vaccessData from './mock-data/vaccess.json'

const delay = (ms = 300) => new Promise(r => setTimeout(r, ms))

export const mockApi = {
  async getAllServicesStatus() {
    await delay(200)
    return JSON.parse(JSON.stringify(servicesData))
  },

  async checkServicesReady() {
    await delay(100)
    return true
  },

  async startServer() { await delay(500); return 'OK' },
  async stopServer() { await delay(500); return 'OK' },
  async startHTTPService() { await delay(300); return 'OK' },
  async stopHTTPService() { await delay(300); return 'OK' },
  async startHTTPSService() { await delay(300); return 'OK' },
  async stopHTTPSService() { await delay(300); return 'OK' },
  async startMySQLService() { await delay(300); return 'OK' },
  async stopMySQLService() { await delay(300); return 'OK' },
  async startPHPService() { await delay(300); return 'OK' },
  async stopPHPService() { await delay(300); return 'OK' },
  async enableProxyService() { await delay(200); return 'OK' },
  async disableProxyService() { await delay(200); return 'OK' },
  async enableACMEService() { await delay(200); return 'OK' },
  async disableACMEService() { await delay(200); return 'OK' },
  async restartAllServices() { await delay(1000); return 'OK' },

  async getSitesList() {
    await delay(200)
    return JSON.parse(JSON.stringify(sitesData))
  },

  async createNewSite(siteJSON) {
    await delay(500)
    return 'OK'
  },

  async deleteSite(host) {
    await delay(500)
    return 'OK'
  },

  async updateSiteCache() {
    await delay(200)
    return 'OK'
  },

  async openSiteFolder(host) {
    await delay(100)
  },

  async uploadCertificate(host, certType, certDataBase64) {
    await delay(300)
    return 'OK'
  },

  async getProxyList() {
    await delay(200)
    return JSON.parse(JSON.stringify(proxiesData))
  },

  async getVAccessRules(host, isProxy) {
    await delay(200)
    return JSON.parse(JSON.stringify(vaccessData))
  },

  async saveVAccessRules(host, isProxy, configJSON) {
    await delay(300)
    return 'OK'
  },

  async getConfig() {
    await delay(200)
    return JSON.parse(JSON.stringify(configData))
  },

  async saveConfig(configJSON) {
    await delay(400)
    return 'OK'
  },

  async reloadConfig() {
    await delay(200)
    return 'OK'
  },

  async obtainSSLCertificate(domain) {
    await delay(1500)
    return 'OK'
  },

  async obtainAllSSLCertificates() {
    await delay(2000)
    return 'OK'
  },

  async getCertInfo(domain) {
    await delay(200)
    const cert = certsData.find(c => c.domain === domain)
    return cert || { has_cert: false }
  },

  async getAllCertsInfo() {
    await delay(300)
    return JSON.parse(JSON.stringify(certsData))
  },

  async deleteCertificate(domain) {
    await delay(300)
    return 'OK'
  },

  async reloadSSLCertificates() {
    await delay(300)
    return 'OK'
  },
}
