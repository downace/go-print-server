import { GetConfig, UpdateServerHost, UpdateServerPort } from "@/go/main/App";
import { defineStore } from "pinia";
import { onBeforeMount, readonly, shallowRef } from "vue";

export const useConfigStore = defineStore("config", () => {
  const isLoaded = shallowRef(false);

  const host = shallowRef("");
  const port = shallowRef(0);

  async function loadConfig() {
    const config = await GetConfig();

    host.value = config.host;
    port.value = config.port;

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

  onBeforeMount(loadConfig);

  return {
    isLoaded: readonly(isLoaded),

    host: readonly(host),
    port: readonly(port),

    updateHost,
    updatePort,
  };
});
