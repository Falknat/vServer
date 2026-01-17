package sites

type SiteInfo struct {
	Name              string   `json:"name"`
	Host              string   `json:"host"`
	Alias             []string `json:"alias"`
	Status            string   `json:"status"`
	RootFile          string   `json:"root_file"`
	RootFileRouting   bool     `json:"root_file_routing"`
	AutoCreateSSL     bool     `json:"auto_create_ssl"`
}

