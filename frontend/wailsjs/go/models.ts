export namespace appconfig {
	
	export interface AuthConfig {
	    enabled: boolean;
	    username: string;
	    password: string;
	}
	export interface TLSConfig {
	    enabled: boolean;
	    certFile: string;
	    keyFile: string;
	}
	export interface AppConfig {
	    host: string;
	    port: number;
	    responseHeaders: Record<string, string>;
	    tls: TLSConfig;
	    auth: AuthConfig;
	}
	

}

export namespace gui {
	
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

