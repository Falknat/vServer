package sites

import (
	config "vServer/Backend/config"
)

func GetSitesList() []SiteInfo {
	sites := make([]SiteInfo, 0)

	for _, site := range config.ConfigData.Site_www {
		siteInfo := SiteInfo{
			Name:            site.Name,
			Host:            site.Host,
			Alias:           site.Alias,
			Status:          site.Status,
			RootFile:        site.Root_file,
			RootFileRouting: site.Root_file_routing,
			AutoCreateSSL:   site.AutoCreateSSL,
		}
		sites = append(sites, siteInfo)
	}

	return sites
}

