import { useConfigStore } from "@/configStore";
import { GetServerStatus, StartServer, StopServer } from "@/go/gui/App";
import type { gui } from "@/go/models";
import { EventsOn } from "@/runtime";
import { defineStore } from "pinia";
import { equals } from "ramda";
import { computed, readonly, shallowRef, watch } from "vue";

export const useServerStore = defineStore("server", () => {
  const configStore = useConfigStore();

  const status = shallowRef<gui.ServerStatus>({
    running: false,
    runningHost: "",
    runningPort: 0,
    error: "",
  });
  const starting = shallowRef(false);
  const error = shallowRef("");

  async function startServer() {
    if (starting.value) {
      return;
    }
    error.value = "";
    starting.value = true;
    try {
      await StartServer();
    } catch (e) {
      error.value = e as string;
    } finally {
      starting.value = false;
    }
  }

  async function stopServer() {
    if (starting.value) {
      return;
    }
    error.value = "";
    starting.value = true;
    try {
      await StopServer();
    } catch (e) {
      error.value = e as string;
    } finally {
      starting.value = false;
    }
  }

  async function toggleServer() {
    return status.value.running ? stopServer() : startServer();
  }

  EventsOn("server-status-changed", (s: gui.ServerStatus) => {
    if (s.running !== status.value.running) {
      needsRestart.value = false;
    }
    status.value = s;
  });

  const needsRestart = shallowRef(false);

  watch(
    () => configStore.config,
    (newConf, oldConf) => {
      if (status.value.running && !equals(newConf, oldConf)) {
        needsRestart.value = true;
      }
    },
  );

  const runningAddr = computed(() =>
    status.value.running
      ? `${status.value.runningHost}:${status.value.runningPort}`
      : null,
  );

  async function refreshStatus() {
    status.value = await GetServerStatus();
  }

  return {
    starting: readonly(starting),
    status: readonly(status),
    error: readonly(error),

    needsRestart: readonly(needsRestart),

    runningAddr,

    refreshStatus,
    toggleServer,
  };
});
