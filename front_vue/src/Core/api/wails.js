const app = () => window.go.admin.App

const methods = [
  'GetAllServicesStatus', 'CheckServicesReady',
  'StartServer', 'StopServer',
  'StartHTTPService', 'StopHTTPService',
  'StartHTTPSService', 'StopHTTPSService',
  'StartMySQLService', 'StopMySQLService',
  'StartPHPService', 'StopPHPService',
  'EnableProxyService', 'DisableProxyService',
  'EnableACMEService', 'DisableACMEService',
  'RestartAllServices',
  'GetSitesList', 'CreateNewSite', 'DeleteSite', 'UpdateSiteCache', 'OpenSiteFolder',
  'UploadCertificate',
  'GetProxyList',
  'GetVAccessRules', 'SaveVAccessRules',
  'GetConfig', 'SaveConfig', 'ReloadConfig',
  'ObtainSSLCertificate', 'ObtainAllSSLCertificates',
  'GetCertInfo', 'GetAllCertsInfo', 'DeleteCertificate', 'ReloadSSLCertificates',
]

const toCamelCase = (s) => s.charAt(0).toLowerCase() + s.slice(1)

export const wailsApi = Object.fromEntries(
  methods.map(m => [toCamelCase(m), (...args) => app()[m](...args)])
)
