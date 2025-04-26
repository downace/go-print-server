<script setup lang="ts">
import AppSettingsHeaders from "@/components/AppSettings/AppSettingsHeaders.vue";
import AppSettingsListen from "@/components/AppSettings/AppSettingsListen.vue";
import AppSettingsTls from "@/components/AppSettings/AppSettingsTls.vue";
import { useConfigStore } from "@/configStore";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";

const configStore = useConfigStore();
</script>

<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="full-height column no-wrap">
      <q-toolbar class="bg-primary text-white">
        <q-toolbar-title> Server Settings </q-toolbar-title>

        <q-space />

        <q-btn flat round icon="mdi-close" to="/" title="Close" />
      </q-toolbar>

      <div class="col column no-wrap q-pa-xs q-gutter-y-sm">
        <template v-if="configStore.isLoaded">
          <app-settings-listen />
          <app-settings-headers />
          <app-settings-tls />
        </template>
        <template v-else>
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
        </template>
      </div>
    </div>
  </q-page>
</template>
