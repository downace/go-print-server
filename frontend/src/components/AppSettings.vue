<script setup lang="ts">
import { useConfigStore } from "@/configStore";
import { GetAvailableAddrs } from "@/go/main/App";
import { main } from "@/go/models";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
import { useDebounce } from "@vueuse/core";
import { isIPv4, isIPv6 } from "is-ip";
import {
  computed,
  onBeforeMount,
  onBeforeUnmount,
  shallowRef,
  watch,
} from "vue";

type IpFamily = "all" | "ipv4" | "ipv6";

const configStore = useConfigStore();

const host = shallowRef<string>("");
const port = shallowRef<number>(0);

watch(
  () => configStore.isLoaded,
  (loaded) => {
    if (loaded) {
      host.value = configStore.host;
      port.value = configStore.port;
    }
  },
  { immediate: true },
);

const showOnlyUp = shallowRef(true);
const showOnlyIpFamily = shallowRef<IpFamily>("ipv4");
const ipFamilies: IpFamily[] = ["all", "ipv4", "ipv6"];

function ipFamilyLabel(family: IpFamily) {
  return {
    all: "All",
    ipv4: "IPv4",
    ipv6: "IPv6",
  }[family];
}

const availableAddrs = shallowRef([] as main.NetInterfaceAddress[]);

const ipsToShow = computed(() =>
  [
    {
      ip: "0.0.0.0",
      interface: {
        name: "All interfaces",
        isUp: true,
      },
    },
    {
      ip: "::",
      interface: {
        name: "All interfaces",
        isUp: true,
      },
    },
    ...availableAddrs.value,
  ].filter(function (addr) {
    if (showOnlyUp.value && !addr.interface.isUp) {
      return false;
    }
    return (
      showOnlyIpFamily.value === "all" ||
      (showOnlyIpFamily.value === "ipv4" && isIPv4(addr.ip)) ||
      (showOnlyIpFamily.value === "ipv6" && isIPv6(addr.ip))
    );
  }),
);

onBeforeMount(async () => {
  availableAddrs.value = await GetAvailableAddrs();
});

const hostToSave = useDebounce(host, 1000);

watch(hostToSave, (h) => {
  configStore.updateHost(h);
});

const portToSave = useDebounce(port, 1000);

watch(portToSave, (p) => {
  configStore.updatePort(p);
});

onBeforeUnmount(async () => {
  await configStore.updateHost(host.value);
  await configStore.updatePort(port.value);
});
</script>

<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="full-height column no-wrap">
      <q-toolbar class="bg-primary text-white">
        <q-toolbar-title> Server Settings </q-toolbar-title>

        <q-space />

        <q-btn flat round icon="mdi-close" to="/" title="Close" />
      </q-toolbar>

      <div class="col relative-position">
        <div class="q-pa-sm">
          <div class="row no-wrap">
            <div class="col">
              <q-select
                v-model="host"
                :options="ipsToShow"
                label="IP address"
                option-value="ip"
                option-label="ip"
                emit-value
              >
                <template #before-options>
                  <div class="row no-wrap items-center">
                    <div class="col">
                      <q-checkbox v-model="showOnlyUp" label="Only UP" />
                    </div>

                    <div class="col">
                      <q-select
                        v-model="showOnlyIpFamily"
                        :options="ipFamilies"
                        label="IP families"
                        :option-label="ipFamilyLabel"
                      >
                      </q-select>
                    </div>
                  </div>
                  <q-separator />
                </template>
                <template #option="{ opt, itemProps }">
                  <q-item v-bind="itemProps">
                    <q-item-section>
                      <q-item-label>
                        {{ opt.ip }}
                      </q-item-label>
                      <q-item-label caption>
                        {{ opt.interface.name }}
                      </q-item-label>
                    </q-item-section>
                    <q-item-section
                      v-if="!opt.interface.isUp"
                      avatar
                      title="Interface is DOWN"
                    >
                      <q-icon name="mdi-link-off" />
                    </q-item-section>
                  </q-item>
                </template>
              </q-select>
            </div>
            <div class="col-3">
              <q-input v-model.number="port" label="Port"></q-input>
            </div>
          </div>
        </div>
        <q-inner-loading :showing="!configStore.isLoaded" />
      </div>
    </div>
  </q-page>
</template>
