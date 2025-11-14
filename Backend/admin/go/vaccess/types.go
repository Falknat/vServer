package vaccess

type VAccessRule struct {
	Type          string   `json:"type"`
	TypeFile      []string `json:"type_file"`
	PathAccess    []string `json:"path_access"`
	IPList        []string `json:"ip_list"`
	ExceptionsDir []string `json:"exceptions_dir"`
	UrlError      string   `json:"url_error"`
}

type VAccessConfig struct {
	Rules []VAccessRule `json:"rules"`
}
