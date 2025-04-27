<script setup lang="ts">
import FilePicker from "@/components/FilePicker.vue";
import useSettingsField from "@/composables/useSettingsField";
import { useConfigStore } from "@/configStore";

const configStore = useConfigStore();

const { value: authEnabled, error: authEnabledError } = useSettingsField(
  () => configStore.config.auth.enabled,
  configStore.updateAuthEnabled,
);

const { value: username, error: usernameError } = useSettingsField(
  () => configStore.config.auth.username,
  configStore.updateAuthUsername,
);

const { value: password, error: passwordError } = useSettingsField(
  () => configStore.config.auth.password,
  configStore.updateAuthPassword,
);
</script>

<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">Authorization</div>
      <q-checkbox v-model="authEnabled" label="Enabled" />
      <div
        v-if="authEnabledError"
        class="text-red ellipsis"
        :title="authEnabledError"
      >
        {{ authEnabledError }}
      </div>
      <div class="row no-wrap q-gutter-x-sm">
        <q-input
          v-if="authEnabled"
          label="Username"
          :model-value="username"
          @change="username = $event"
          @input="usernameError = ''"
          :error="!!usernameError"
          :error-message="usernameError"
        />
        <q-input
          v-if="authEnabled"
          label="Password"
          :model-value="password"
          @change="password = $event"
          @input="passwordError = ''"
          :error="!!passwordError"
          :error-message="passwordError"
        />
      </div>
    </q-card-section>
  </q-card>
</template>
