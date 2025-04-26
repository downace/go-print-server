<script setup lang="ts">
import { GetAvailableAddrs } from "@/go/gui/App";
import { gui } from "@/go/models";
import { isIPv4, isIPv6 } from "is-ip";
import { computed, onBeforeMount, shallowRef } from "vue";

type IpFamily = "all" | "ipv4" | "ipv6";

const { error } = defineProps<{
  error: string;
}>();

const host = defineModel();

const showOnlyUp = shallowRef(true);
const showOnlyIpFamily = shallowRef<IpFamily>("ipv4");
const ipFamilies: IpFamily[] = ["all", "ipv4", "ipv6"];

const availableAddrs = shallowRef([] as gui.NetInterfaceAddress[]);

onBeforeMount(async () => {
  availableAddrs.value = await GetAvailableAddrs();
});

function ipFamilyLabel(family: IpFamily) {
  return {
    all: "All",
    ipv4: "IPv4",
    ipv6: "IPv6",
  }[family];
}

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
</script>

<template>
  <q-select
    v-model="host"
    :options="ipsToShow"
    label="IP"
    option-value="ip"
    option-label="ip"
    emit-value
    :error="!!error"
    :error-message="error"
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
</template>
