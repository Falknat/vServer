package proxy

type ProxyInfo struct {
	Name            string `json:"Name"`
	Enable          bool   `json:"Enable"`
	ExternalDomain  string `json:"ExternalDomain"`
	LocalAddress    string `json:"LocalAddress"`
	LocalPort       string `json:"LocalPort"`
	ServiceHTTPSuse bool   `json:"ServiceHTTPSuse"`
	AutoHTTPS       bool   `json:"AutoHTTPS"`
	AutoCreateSSL   bool   `json:"AutoCreateSSL"`
	Compression     *bool  `json:"Compression"`
	Status          string `json:"Status"`
}

