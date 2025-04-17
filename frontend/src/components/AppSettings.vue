<script setup lang="ts">
import { useConfigStore } from "@/configStore";
import { GetAvailableAddrs, PickFilePath } from "@/go/gui/App";
import { gui } from "@/go/models";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
import { isIPv4, isIPv6 } from "is-ip";
import { computed, onBeforeMount, shallowRef, watch } from "vue";

type IpFamily = "all" | "ipv4" | "ipv6";

const configStore = useConfigStore();

const host = shallowRef("");
const port = shallowRef("0");
const responseHeadersStr = shallowRef("");
const tlsEnabled = shallowRef(false);
const tlsCertFile = shallowRef("");
const tlsKeyFile = shallowRef("");

watch(
  () => configStore.isLoaded,
  (loaded) => {
    if (loaded) {
      host.value = configStore.host;
      port.value = String(configStore.port);
      tlsEnabled.value = configStore.tlsEnabled;
      tlsCertFile.value = configStore.tlsCertFile;
      tlsKeyFile.value = configStore.tlsKeyFile;
      responseHeadersStr.value = [...configStore.responseHeaders.entries()]
        .map(([name, value]) => `${name}: ${value}`)
        .join("\n");
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

const availableAddrs = shallowRef([] as gui.NetInterfaceAddress[]);

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

const hostError = shallowRef("");
const portError = shallowRef("");
const headersError = shallowRef("");
const tlsEnabledError = shallowRef("");
const tlsCertFileError = shallowRef("");
const tlsKeyFileError = shallowRef("");

watch(host, async (host) => {
  try {
    await configStore.updateHost(host);
  } catch (e) {
    hostError.value = e as string;
  }
});

watch(port, async (port) => {
  try {
    const portNumber = Number.parseInt(port);
    if (!Number.isFinite(portNumber)) {
      portError.value = "Invalid port number";
      return;
    }
    await configStore.updatePort(portNumber);
  } catch (e) {
    portError.value = e as string;
  }
});

watch(responseHeadersStr, async (headers) => {
  try {
    await configStore.updateResponseHeaders(
      new Map(
        headers
          .split("\n")
          .filter((line) => line.trim() !== "")
          .map((line) => {
            const pair = line.split(":", 2);
            return [pair[0].trim(), (pair[1] ?? "").trim()];
          }),
      ),
    );
  } catch (e) {
    headersError.value = e as string;
  }
});

watch(tlsEnabled, async (enabled) => {
  try {
    await configStore.updateTlsEnabled(enabled);
  } catch (e) {
    tlsEnabledError.value = e as string;
  }
});

watch(tlsCertFile, async (certFile) => {
  try {
    await configStore.updateTlsCertFile(certFile);
  } catch (e) {
    tlsCertFileError.value = e as string;
  }
});

watch(tlsKeyFile, async (keyFile) => {
  try {
    await configStore.updateTlsKeyFile(keyFile);
  } catch (e) {
    tlsKeyFileError.value = e as string;
  }
});

async function pickTlsCertFile() {
  const newPath = await PickFilePath();

  if (!newPath) {
    return;
  }

  tlsCertFile.value = newPath;
}

async function pickTlsKeyFile() {
  const newPath = await PickFilePath();

  if (!newPath) {
    return;
  }

  tlsKeyFile.value = newPath;
}

const responseHeadersPlaceholder = `Example:

Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: Accept
`;
</script>

<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="full-height column no-wrap">
      <q-toolbar class="bg-primary text-white">
        <q-toolbar-title> Server Settings </q-toolbar-title>

        <q-space />

        <q-btn flat round icon="mdi-close" to="/" title="Close" />
      </q-toolbar>

      <div class="col relative-position column no-wrap q-pa-xs q-gutter-y-sm">
        <q-card>
          <q-card-section>
            <div class="text-h6">Server Address</div>
            <div class="row no-wrap">
              <div class="col">
                <q-select
                  v-model="host"
                  :options="ipsToShow"
                  label="IP"
                  option-value="ip"
                  option-label="ip"
                  emit-value
                  :error="!!hostError"
                  :error-message="hostError"
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
        <q-card>
          <q-card-section>
            <div class="text-h6">Response Headers</div>
            <q-input
              type="textarea"
              rows="4"
              :model-value="responseHeadersStr"
              @change="responseHeadersStr = $event"
              @input="headersError = ''"
              :placeholder="responseHeadersPlaceholder"
              :error="!!headersError"
              :error-message="headersError"
            />
          </q-card-section>
        </q-card>
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
              <q-input
                v-if="tlsEnabled"
                label="Cert File"
                :model-value="tlsCertFile"
                @change="tlsCertFile = $event"
                @input="tlsCertFileError = ''"
                :error="!!tlsCertFileError"
                :error-message="tlsCertFileError"
              >
                <template #append>
                  <q-btn
                    flat
                    round
                    icon="mdi-file"
                    title="Pick File"
                    @click="pickTlsCertFile"
                  />
                </template>
              </q-input>
              <q-input
                v-if="tlsEnabled"
                label="Key File"
                :model-value="tlsKeyFile"
                @change="tlsKeyFile = $event"
                @input="tlsKeyFileError = ''"
                :error="!!tlsKeyFileError"
                :error-message="tlsKeyFileError"
              >
                <template #append>
                  <q-btn
                    flat
                    round
                    icon="mdi-file"
                    title="Pick File"
                    @click="pickTlsKeyFile"
                  />
                </template>
              </q-input>
            </div>
          </q-card-section>
        </q-card>
        <q-inner-loading :showing="!configStore.isLoaded" />
      </div>
    </div>
  </q-page>
</template>
