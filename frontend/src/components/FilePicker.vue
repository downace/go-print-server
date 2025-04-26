<script setup lang="ts">
import { PickFilePath } from "@/go/gui/App";

const file = defineModel<string>({ required: true });
const error = defineModel<string>("error", { default: "" });

async function pickFile() {
  const newPath = await PickFilePath();

  if (!newPath) {
    return;
  }

  file.value = newPath;
}
</script>

<template>
  <q-input
    :model-value="file"
    @change="file = $event"
    @input="error = ''"
    :error="!!error"
    :error-message="error"
  >
    <template #append>
      <q-btn flat round icon="mdi-file" title="Pick File" @click="pickFile" />
    </template>
  </q-input>
</template>
