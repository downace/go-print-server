<script setup lang="ts">
import ListenHostPicker from "@/components/ListenHostPicker.vue";
import useSettingsField from "@/composables/useSettingsField";
import { useConfigStore } from "@/configStore";

const configStore = useConfigStore();

const { value: host, error: hostError } = useSettingsField(
  () => configStore.config.host,
  configStore.updateHost,
);

const { value: port, error: portError } = useSettingsField(
  () => String(configStore.config.port),
  (newPort) => {
    const portNumber = Number.parseInt(newPort);
    if (!Number.isFinite(portNumber)) {
      throw new Error("Invalid port number");
    }
    return configStore.updatePort(portNumber);
  },
);
</script>

<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">Server Address</div>
      <div class="row no-wrap">
        <div class="col">
          <listen-host-picker v-model="host" :error="hostError" />
        </div>
        <div class="col-3">
          <q-input
            :model-value="port"
            @change="port = $event"
            @input="portError = ''"
            type="number"
            label="Port"
            :error="!!portError"
            :error-message="portError"
          />
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>
