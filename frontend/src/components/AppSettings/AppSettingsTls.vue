<script setup lang="ts">
import FilePicker from "@/components/FilePicker.vue";
import useSettingsField from "@/composables/useSettingsField";
import { useConfigStore } from "@/configStore";

const configStore = useConfigStore();

const { value: tlsEnabled, error: tlsEnabledError } = useSettingsField(
  () => configStore.config.tls.enabled,
  configStore.updateTlsEnabled,
);

const {
  value: certFile,
  error: certFileError,
  clearError: clearCertFileError,
} = useSettingsField(
  () => configStore.config.tls.certFile,
  configStore.updateTlsCertFile,
);

const {
  value: keyFile,
  error: keyFileError,
  clearError: clearKeyFileError,
} = useSettingsField(
  () => configStore.config.tls.keyFile,
  configStore.updateTlsKeyFile,
);
</script>

<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">TLS</div>
      <q-checkbox v-model="tlsEnabled" label="Enabled" />
      <div
        v-if="tlsEnabledError"
        class="text-red ellipsis"
        :title="tlsEnabledError"
      >
        {{ tlsEnabledError }}
      </div>
      <div class="row no-wrap q-gutter-x-sm">
        <file-picker
          v-if="tlsEnabled"
          label="Cert File"
          v-model="certFile"
          :error="certFileError"
          @update:error="clearCertFileError()"
        />
        <file-picker
          v-if="tlsEnabled"
          label="Key File"
          v-model="keyFile"
          :error="keyFileError"
          @update:error="clearKeyFileError()"
        />
      </div>
    </q-card-section>
  </q-card>
</template>
