package proxy

type ProxyInfo struct {
	Enable          bool   `json:"enable"`
	ExternalDomain  string `json:"external_domain"`
	LocalAddress    string `json:"local_address"`
	LocalPort       string `json:"local_port"`
	ServiceHTTPSuse bool   `json:"service_https_use"`
	AutoHTTPS       bool   `json:"auto_https"`
	AutoCreateSSL   bool   `json:"auto_create_ssl"`
	Status          string `json:"status"`
}

