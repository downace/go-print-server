<script setup lang="ts">
import useSettingsField from "@/composables/useSettingsField";
import { useConfigStore } from "@/configStore";

const configStore = useConfigStore();

const { value: headersStr, error: headersError } = useSettingsField(
  () =>
    Object.entries(configStore.config.responseHeaders)
      .map(([name, value]) => `${name}: ${value}`)
      .join("\n"),
  (headers) =>
    configStore.updateResponseHeaders(
      Object.fromEntries(
        headers
          .split("\n")
          .filter((line) => line.trim() !== "")
          .map((line) => {
            const pair = line.split(":", 2);
            return [pair[0].trim(), (pair[1] ?? "").trim()];
          }),
      ),
    ),
);

const responseHeadersPlaceholder = `Example:

Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: Accept
`;
</script>

<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">Response Headers</div>
      <q-input
        type="textarea"
        rows="4"
        :model-value="headersStr"
        @change="headersStr = $event"
        @input="headersError = ''"
        :placeholder="responseHeadersPlaceholder"
        :error="!!headersError"
        :error-message="headersError"
      />
    </q-card-section>
  </q-card>
</template>
