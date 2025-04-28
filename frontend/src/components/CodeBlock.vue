<script setup lang="ts">
import HighlightJsVue from "@highlightjs/vue-plugin";
import "highlight.js/styles/default.css";
import hljs from "highlight.js/lib/core";
import langGo from "highlight.js/lib/languages/go";
import langJs from "highlight.js/lib/languages/javascript";
import langShell from "highlight.js/lib/languages/shell";
import { copyToClipboard, QBtn } from "quasar";

const HighlightJs = HighlightJsVue.component;

hljs.registerLanguage("js", langJs);
hljs.registerLanguage("go", langGo);
hljs.registerLanguage("shell", langShell);

const { language, code } = defineProps<{
  label: string;
  language: string;
  code: string;
}>();
</script>

<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">{{ label }}</div>

      <q-btn
        class="absolute-top-right q-ma-sm"
        flat
        dense
        icon="mdi-content-copy"
        @click="copyToClipboard(code)"
      >
        <q-tooltip> Copy to clipboard </q-tooltip>
      </q-btn>
    </q-card-section>
    <highlight-js :language="language" :code="code" class="q-ma-none" />
  </q-card>
</template>
