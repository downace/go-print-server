import {
  GetConfig,
  UpdateAuthEnabled,
  UpdateAuthPassword,
  UpdateAuthUsername,
  UpdateResponseHeaders,
  UpdateServerHost,
  UpdateServerPort,
  UpdateTLSCertFile,
  UpdateTLSEnabled,
  UpdateTLSKeyFile,
} from "@/go/gui/App";
import { appconfig } from "@/go/models";
import { defineStore } from "pinia";
import { equals, mergeDeepRight } from "ramda";
import { PartialDeep } from "type-fest";
import { onBeforeMount, readonly, shallowRef } from "vue";

export const useConfigStore = defineStore("config", () => {
  const isLoaded = shallowRef(false);

  const config = shallowRef<appconfig.AppConfig>({
    host: "",
    port: 0,
    responseHeaders: {},
    tls: {
      enabled: false,
      certFile: "",
      keyFile: "",
    },
    auth: {
      enabled: false,
      username: "",
      password: "",
    },
  });

  async function loadConfig() {
    config.value = Object.freeze(await GetConfig());

    isLoaded.value = true;
  }

  async function updateConfig(
    update: () => Promise<unknown>,
    localPatch: PartialDeep<appconfig.AppConfig>,
  ) {
    const newConfig = mergeDeepRight(config.value, localPatch);
    if (equals(newConfig, config.value)) {
      return;
    }
    await update();
    config.value = Object.freeze(newConfig);
  }

  async function updateHost(newHost: string) {
    await updateConfig(() => UpdateServerHost(newHost), { host: newHost });
  }

  async function updatePort(newPort: number) {
    await updateConfig(() => UpdateServerPort(newPort), { port: newPort });
  }

  async function updateResponseHeaders(newHeaders: Record<string, string>) {
    await updateConfig(() => UpdateResponseHeaders(newHeaders), {
      responseHeaders: newHeaders,
    });
  }

  async function updateTlsEnabled(enabled: boolean) {
    await updateConfig(() => UpdateTLSEnabled(enabled), { tls: { enabled } });
  }

  async function updateTlsCertFile(certFile: string) {
    await updateConfig(() => UpdateTLSCertFile(certFile), {
      tls: { certFile },
    });
  }

  async function updateTlsKeyFile(keyFile: string) {
    await updateConfig(() => UpdateTLSKeyFile(keyFile), { tls: { keyFile } });
  }

  async function updateAuthEnabled(enabled: boolean) {
    await updateConfig(() => UpdateAuthEnabled(enabled), { auth: { enabled } });
  }

  async function updateAuthUsername(username: string) {
    await updateConfig(() => UpdateAuthUsername(username), {
      auth: { username },
    });
  }

  async function updateAuthPassword(password: string) {
    await updateConfig(() => UpdateAuthPassword(password), {
      auth: { password },
    });
  }

  onBeforeMount(loadConfig);

  return {
    isLoaded: readonly(isLoaded),

    config: readonly(config),

    updateHost,
    updatePort,
    updateResponseHeaders,
    updateTlsEnabled,
    updateTlsCertFile,
    updateTlsKeyFile,
    updateAuthEnabled,
    updateAuthUsername,
    updateAuthPassword,
  };
});
