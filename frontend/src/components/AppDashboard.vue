<script lang="ts" setup>
import { useConfigStore } from "@/configStore";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
import { useServerStore } from "@/serverStore";
import { copyToClipboard } from "quasar";
import { computed } from "vue";

const serverStore = useServerStore();

const configStore = useConfigStore();

const statusToggle = computed({
  get: () => {
    if (serverStore.starting) {
      return null;
    }

    return serverStore.status.running;
  },
  set() {
    serverStore.toggleServer();
  },
});

const serverAddress = computed(
  () =>
    serverStore.runningAddr ??
    `${configStore.config.host}:${configStore.config.port}`,
);
</script>

<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="full-height column no-wrap">
      <q-toolbar>
        <q-btn
          flat
          round
          icon="mdi-file-code"
          to="/snippets"
          title="Code snippets"
        />

        <q-space />

        <q-btn flat round icon="mdi-cog" to="/settings" title="Settings" />
      </q-toolbar>

      <div class="col column no-wrap justify-center items-center q-gutter-lg">
        <q-chip
          size="lg"
          square
          icon-right="mdi-content-copy"
          title="Click to copy"
          clickable
          @click="copyToClipboard(serverAddress)"
        >
          {{ serverAddress }}
        </q-chip>

        <q-toggle
          v-model="statusToggle"
          size="100px"
          dense
          checked-icon="mdi-check"
          color="positive"
          :title="serverStore.status.running ? 'Stop' : 'Start'"
          :disable="serverStore.starting"
        />
        <div v-if="serverStore.status.running">Click to stop server</div>
        <div v-else>Click to start server</div>
      </div>

      <div class="col-2 column justify-end">
        <q-banner v-if="serverStore.needsRestart" class="bg-blue text-white">
          <template #avatar>
            <q-icon name="mdi-restart"></q-icon>
          </template>
          Settings changed, restart server to apply them
        </q-banner>
        <q-banner
          v-if="serverStore.status.error"
          dense
          class="full-width bg-red text-white"
        >
          <template #avatar> </template>
          <div class="ellipsis-3-lines" :title="serverStore.status.error">
            {{ serverStore.status.error }}
          </div>
        </q-banner>
      </div>
    </div>
  </q-page>
</template>
