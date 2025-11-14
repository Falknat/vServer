package services

type ServiceStatus struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Port   string `json:"port"`
	Info   string `json:"info"`
}

type AllServicesStatus struct {
	HTTP  ServiceStatus `json:"http"`
	HTTPS ServiceStatus `json:"https"`
	MySQL ServiceStatus `json:"mysql"`
	PHP   ServiceStatus `json:"php"`
	Proxy ServiceStatus `json:"proxy"`
}
