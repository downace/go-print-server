import {
  GetConfig,
  UpdateResponseHeaders,
  UpdateServerHost,
  UpdateServerPort,
} from "@/go/main/App";
import { defineStore } from "pinia";
import { equals } from "ramda";
import { onBeforeMount, readonly, shallowRef } from "vue";

export const useConfigStore = defineStore("config", () => {
  const isLoaded = shallowRef(false);

  const host = shallowRef("");
  const port = shallowRef(0);
  const responseHeaders = shallowRef(new Map<string, string>());

  async function loadConfig() {
    const config = await GetConfig();

    host.value = config.host;
    port.value = config.port;
    responseHeaders.value = new Map(Object.entries(config.responseHeaders));

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

  onBeforeMount(loadConfig);

  return {
    isLoaded: readonly(isLoaded),

    host: readonly(host),
    port: readonly(port),
    responseHeaders: readonly(responseHeaders),

    updateHost,
    updatePort,
    updateResponseHeaders,
  };
});
