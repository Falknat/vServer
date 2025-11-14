package proxy

import (
	config "vServer/Backend/config"
)

func GetProxyList() []ProxyInfo {
	proxies := make([]ProxyInfo, 0)

	for _, proxyConfig := range config.ConfigData.Proxy_Service {
		status := "disabled"
		if proxyConfig.Enable {
			status = "active"
		}

		proxyInfo := ProxyInfo{
			Enable:          proxyConfig.Enable,
			ExternalDomain:  proxyConfig.ExternalDomain,
			LocalAddress:    proxyConfig.LocalAddress,
			LocalPort:       proxyConfig.LocalPort,
			ServiceHTTPSuse: proxyConfig.ServiceHTTPSuse,
			AutoHTTPS:       proxyConfig.AutoHTTPS,
			Status:          status,
		}
		proxies = append(proxies, proxyInfo)
	}

	return proxies
}

