const app = () => window.go.admin.App

export const wailsApi = {
  async getAllServicesStatus() {
    return await app().GetAllServicesStatus()
  },

  async checkServicesReady() {
    return await app().CheckServicesReady()
  },

  async startServer() { return await app().StartServer() },
  async stopServer() { return await app().StopServer() },
  async startHTTPService() { return await app().StartHTTPService() },
  async stopHTTPService() { return await app().StopHTTPService() },
  async startHTTPSService() { return await app().StartHTTPSService() },
  async stopHTTPSService() { return await app().StopHTTPSService() },
  async startMySQLService() { return await app().StartMySQLService() },
  async stopMySQLService() { return await app().StopMySQLService() },
  async startPHPService() { return await app().StartPHPService() },
  async stopPHPService() { return await app().StopPHPService() },
  async enableProxyService() { return await app().EnableProxyService() },
  async disableProxyService() { return await app().DisableProxyService() },
  async enableACMEService() { return await app().EnableACMEService() },
  async disableACMEService() { return await app().DisableACMEService() },
  async restartAllServices() { return await app().RestartAllServices() },

  async getSitesList() {
    return await app().GetSitesList()
  },

  async createNewSite(siteJSON) {
    return await app().CreateNewSite(siteJSON)
  },

  async deleteSite(host) {
    return await app().DeleteSite(host)
  },

  async updateSiteCache() {
    return await app().UpdateSiteCache()
  },

  async openSiteFolder(host) {
    return await app().OpenSiteFolder(host)
  },

  async uploadCertificate(host, certType, certDataBase64) {
    return await app().UploadCertificate(host, certType, certDataBase64)
  },

  async getProxyList() {
    return await app().GetProxyList()
  },

  async getVAccessRules(host, isProxy) {
    return await app().GetVAccessRules(host, isProxy)
  },

  async saveVAccessRules(host, isProxy, configJSON) {
    return await app().SaveVAccessRules(host, isProxy, configJSON)
  },

  async getConfig() {
    return await app().GetConfig()
  },

  async saveConfig(configJSON) {
    return await app().SaveConfig(configJSON)
  },

  async reloadConfig() {
    return await app().ReloadConfig()
  },

  async obtainSSLCertificate(domain) {
    return await app().ObtainSSLCertificate(domain)
  },

  async obtainAllSSLCertificates() {
    return await app().ObtainAllSSLCertificates()
  },

  async getCertInfo(domain) {
    return await app().GetCertInfo(domain)
  },

  async getAllCertsInfo() {
    return await app().GetAllCertsInfo()
  },

  async deleteCertificate(domain) {
    return await app().DeleteCertificate(domain)
  },

  async reloadSSLCertificates() {
    return await app().ReloadSSLCertificates()
  },
}
