import {
  GetConfig,
  UpdateResponseHeaders,
  UpdateServerHost,
  UpdateServerPort,
  UpdateTLSCertFile,
  UpdateTLSEnabled,
  UpdateTLSKeyFile,
} from "@/go/gui/App";
import { defineStore } from "pinia";
import { equals } from "ramda";
import { onBeforeMount, readonly, shallowRef } from "vue";

export const useConfigStore = defineStore("config", () => {
  const isLoaded = shallowRef(false);

  const host = shallowRef("");
  const port = shallowRef(0);
  const responseHeaders = shallowRef(new Map<string, string>());
  const tlsEnabled = shallowRef(false);
  const tlsCertFile = shallowRef("");
  const tlsKeyFile = shallowRef("");

  async function loadConfig() {
    const config = await GetConfig();

    host.value = config.host;
    port.value = config.port;
    responseHeaders.value = new Map(Object.entries(config.responseHeaders));
    tlsEnabled.value = config.tls.enabled;
    tlsCertFile.value = config.tls.certFile;
    tlsKeyFile.value = config.tls.keyFile;

    isLoaded.value = true;
  }

  async function updateHost(newHost: string) {
    if (newHost === host.value) {
      return;
    }
    await UpdateServerHost(newHost);
    host.value = newHost;
  }

  async function updatePort(newPort: number) {
    if (newPort === port.value) {
      return;
    }
    await UpdateServerPort(newPort);
    port.value = newPort;
  }

  async function updateResponseHeaders(newHeaders: Map<string, string>) {
    if (equals(newHeaders, responseHeaders.value)) {
      return;
    }
    await UpdateResponseHeaders(Object.fromEntries(newHeaders.entries()));
    responseHeaders.value = newHeaders;
  }

  async function updateTlsEnabled(enabled: boolean) {
    if (enabled === tlsEnabled.value) {
      return;
    }
    await UpdateTLSEnabled(enabled);
    tlsEnabled.value = enabled;
  }

  async function updateTlsCertFile(newFile: string) {
    if (newFile === tlsCertFile.value) {
      return;
    }
    await UpdateTLSCertFile(newFile);
    tlsCertFile.value = newFile;
  }

  async function updateTlsKeyFile(newFile: string) {
    if (newFile === tlsKeyFile.value) {
      return;
    }
    await UpdateTLSKeyFile(newFile);
    tlsKeyFile.value = newFile;
  }

  onBeforeMount(loadConfig);

  return {
    isLoaded: readonly(isLoaded),

    host: readonly(host),
    port: readonly(port),
    tlsEnabled: readonly(tlsEnabled),
    tlsCertFile: readonly(tlsCertFile),
    tlsKeyFile: readonly(tlsKeyFile),
    responseHeaders: readonly(responseHeaders),

    updateHost,
    updatePort,
    updateResponseHeaders,
    updateTlsEnabled,
    updateTlsCertFile,
    updateTlsKeyFile,
  };
});
