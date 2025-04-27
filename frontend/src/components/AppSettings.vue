<script setup lang="ts">
import AppSettingsAuth from "@/components/AppSettings/AppSettingsAuth.vue";
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

      <div class="col column no-wrap q-pa-xs scroll-y">
        <div v-if="configStore.isLoaded" class="q-gutter-y-sm">
          <app-settings-listen />
          <app-settings-headers />
          <app-settings-tls />
          <app-settings-auth />
        </div>
        <div v-else class="full-height q-gutter-y-sm">
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
          <q-skeleton type="rect" height="10%" />
        </div>
      </div>
    </div>
  </q-page>
</template>
