export namespace main {
	
	export interface AppConfig {
	    host: string;
	    port: number;
	    responseHeaders: Record<string, string>;
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

