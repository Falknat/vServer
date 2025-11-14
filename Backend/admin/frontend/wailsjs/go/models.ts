export namespace proxy {
	
	export class ProxyInfo {
	    enable: boolean;
	    external_domain: string;
	    local_address: string;
	    local_port: string;
	    service_https_use: boolean;
	    auto_https: boolean;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.external_domain = source["external_domain"];
	        this.local_address = source["local_address"];
	        this.local_port = source["local_port"];
	        this.service_https_use = source["service_https_use"];
	        this.auto_https = source["auto_https"];
	        this.status = source["status"];
	    }
	}

}

export namespace services {
	
	export class ServiceStatus {
	    name: string;
	    status: boolean;
	    port: string;
	    requests: number;
	    info: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.port = source["port"];
	        this.requests = source["requests"];
	        this.info = source["info"];
	    }
	}
	export class AllServicesStatus {
	    http: ServiceStatus;
	    https: ServiceStatus;
	    mysql: ServiceStatus;
	    php: ServiceStatus;
	    proxy: ServiceStatus;
	
	    static createFrom(source: any = {}) {
	        return new AllServicesStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.http = this.convertValues(source["http"], ServiceStatus);
	        this.https = this.convertValues(source["https"], ServiceStatus);
	        this.mysql = this.convertValues(source["mysql"], ServiceStatus);
	        this.php = this.convertValues(source["php"], ServiceStatus);
	        this.proxy = this.convertValues(source["proxy"], ServiceStatus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace sites {
	
	export class SiteInfo {
	    name: string;
	    host: string;
	    alias: string[];
	    status: string;
	    root_file: string;
	    root_file_routing: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SiteInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.host = source["host"];
	        this.alias = source["alias"];
	        this.status = source["status"];
	        this.root_file = source["root_file"];
	        this.root_file_routing = source["root_file_routing"];
	    }
	}

}

