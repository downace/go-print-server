export namespace main {
	
	export interface AppConfigTLS {
	    enabled: boolean;
	    certFile: string;
	    keyFile: string;
	}
	export interface AppConfig {
	    host: string;
	    port: number;
	    responseHeaders: Record<string, string>;
	    tls: AppConfigTLS;
	}
	
	export interface NetInterface {
	    name: string;
	    isUp: boolean;
	}
	export interface NetInterfaceAddress {
	    ip: string;
	    interface: NetInterface;
	}
	export interface ServerStatus {
	    running: boolean;
	    error: string;
	    runningHost: string;
	    runningPort: number;
	}

}

