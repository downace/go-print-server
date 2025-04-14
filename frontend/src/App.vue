<script lang="ts" setup>
import { useServerStore } from "@/serverStore";
import { onBeforeMount } from "vue";

const serverStore = useServerStore();

onBeforeMount(() => {
  void serverStore.refreshStatus();
});
</script>

<template>
  <q-layout view="hHh lpR fFf">
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>
