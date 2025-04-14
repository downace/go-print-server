import { GetServerStatus, StartServer, StopServer } from "@/go/main/App";
import type { main } from "@/go/models";
import { EventsOn } from "@/runtime";
import { defineStore } from "pinia";
import { computed, onBeforeMount, readonly, shallowRef } from "vue";

export const useServerStore = defineStore("server", () => {
  const status = shallowRef<main.ServerStatus>({
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

  EventsOn("server-status-changed", (s) => {
    status.value = s;
  });

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

    runningAddr,

    refreshStatus,
    toggleServer,
  };
});
