<script setup lang="ts">
import CodeBlock from "@/components/CodeBlock.vue";
import { useConfigStore } from "@/configStore";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
import Handlebars from "handlebars";
import { QBtn } from "quasar";
import { computed, shallowRef } from "vue";

type Snippet = {
  id: string;
  label: string;
  language: string;
  codeBlocks: {
    label: string;
    template: HandlebarsTemplateDelegate;
  }[];
};

Handlebars.registerHelper("encodeURIComponent", (str) =>
  encodeURIComponent(str),
);

const configStore = useConfigStore();

const snippets: Snippet[] = [
  {
    id: "js",
    label: "JavaScript",
    language: "js",
    codeBlocks: [
      {
        label: "Printers list",
        template:
          Handlebars.compile(`fetch('{{protocol}}://{{host}}:{{port}}/printers'{{#if hasHeaders}}, {
  headers: {
{{#each headers}}
    "{{@key}}": "{{this}}",
{{/each}}
  }
}{{/if}})
  .then(res => res.json())
  .then(res => {
    console.log(res.printers)
  })
`),
      },
      {
        label: "Print PDF from file",
        template:
          Handlebars.compile(`import { promises as fsPromises } from 'fs';

fetch('{{protocol}}://{{host}}:{{port}}/print-pdf?printer={{encodeURIComponent samplePrinter}}', {
  method: 'POST',
  body: await fsPromises.readFile("{{samplePdfPath}}"),
  headers: {
    "Content-Type": "application/pdf",
{{#each headers}}
    "{{@key}}": "{{this}}",
{{/each}}
  }
}).then(res => res.json()).then((res) => {
  console.log(res)
})
`),
      },
      {
        label: "Print PDF from URL",
        template:
          Handlebars.compile(`fetch('{{protocol}}://{{host}}:{{port}}/print-pdf-url?printer={{encodeURIComponent samplePrinter}}&url={{encodeURIComponent samplePrintUrl}}', {
  method: 'POST',
{{#if hasHeaders}}
  headers: {
{{#each headers}}
    "{{@key}}": "{{this}}",
{{/each}}
{{/if}}
  }
}).then(res => res.json()).then((res) => {
  console.log(res)
})`),
      },
    ],
  },
  {
    id: "go",
    label: "Go",
    language: "go",
    codeBlocks: [
      {
        label: "Printers list",
        template: Handlebars.compile(`type Printer struct {
	Name string \`json:"name"\`
}

type Response struct {
	Printers []Printer \`json:"printers"\`
}

req, _ := http.NewRequest("GET", "{{protocol}}://{{host}}:{{port}}/printers", nil)
{{#each headers}}
req.Header.Set("{{@key}}", "{{this}}")
{{/each}}

res, _ := http.DefaultClient.Do(req)

fmt.Println(res.Status)

bodyRaw, _ := io.ReadAll(res.Body)

var bodyParsed Response
_ = json.Unmarshal(bodyRaw, &bodyParsed)

fmt.Printf("%#v\\n", bodyParsed.Printers)`),
      },
      {
        label: "Print PDF from file",
        template: Handlebars.compile(`file, _ := os.Open("{{samplePdfPath}}")
req, _ := http.NewRequest("POST", "{{protocol}}://{{host}}:{{port}}/print-pdf?printer={{encodeURIComponent samplePrinter}}", file)
{{#each headers}}
req.Header.Set("{{@key}}", "{{this}}")
{{/each}}
req.Header.Set("Content-Type", "application/pdf")

res, _ := http.DefaultClient.Do(req)

fmt.Println(res.Status)`),
      },
      {
        label: "Print PDF from URL",
        template:
          Handlebars.compile(`req, _ := http.NewRequest("POST", "{{protocol}}://{{host}}:{{port}}/print-pdf-url?printer={{encodeURIComponent samplePrinter}}&url={{encodeURIComponent samplePrintUrl}}", nil)
{{#each headers}}
req.Header.Set("{{@key}}", "{{this}}")
{{/each}}

res, _ := http.DefaultClient.Do(req)

fmt.Println(res.Status)`),
      },
    ],
  },
  {
    id: "curl",
    label: "curl",
    language: "shell",
    codeBlocks: [
      {
        label: "Printers list",
        template: Handlebars.compile(`curl \\
{{#each headers}}
  -H '{{@key}}: {{this}}' \\
{{/each}}
  {{protocol}}://{{host}}:{{port}}/printers
`),
      },
      {
        label: "Print PDF from file",
        template: Handlebars.compile(`curl \\
{{#each headers}}
  -H '{{@key}}: {{this}}' \\
{{/each}}
  -H 'Content-Type: application/pdf' \\
  --data-binary '{{samplePdfPath}}' \\
  {{protocol}}://{{host}}:{{port}}/print-pdf?printer={{encodeURIComponent samplePrinter}}
`),
      },
      {
        label: "Print PDF from URL",
        template: Handlebars.compile(`curl \\
{{#each headers}}
  -H '{{@key}}: {{this}}' \\
{{/each}}
  {{protocol}}://{{host}}:{{port}}/print-pdf-url?printer={{encodeURIComponent samplePrinter}}&url={{encodeURIComponent samplePrintUrl}}
`),
      },
    ],
  },
];

const tab = shallowRef(snippets[0].id);

const samplePdfPath = shallowRef("/path/to/file.pdf");
const samplePrinter = shallowRef("my printer");
const samplePrintUrl = shallowRef("https://pdfobject.com/pdf/sample.pdf");

const snippetTemplateContext = computed(() => ({
  protocol: configStore.config.tls.enabled ? "https" : "http",
  host: configStore.config.host,
  port: configStore.config.port,
  hasHeaders: configStore.config.auth.enabled,
  headers: configStore.config.auth.enabled
    ? {
        Authorization:
          "Basic " +
          btoa(
            configStore.config.auth.username +
              ":" +
              configStore.config.auth.password,
          ),
      }
    : {},
  samplePdfPath: samplePdfPath.value,
  samplePrinter: samplePrinter.value,
  samplePrintUrl: samplePrintUrl.value,
}));
</script>

<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="full-height column">
      <q-tabs v-model="tab" no-caps :breakpoint="0" align="left">
        <q-tab v-for="s in snippets" :name="s.id">
          {{ s.label }}
        </q-tab>
        <q-space />
        <q-btn flat stretch icon="mdi-close" to="/" title="Close" />
      </q-tabs>
      <q-tab-panels v-model="tab" class="col scroll-y full-width">
        <q-tab-panel v-for="s in snippets" :name="s.id">
          <div class="column no-wrap q-gutter-y-md">
            <q-card>
              <q-card-section>
                <div class="text-h6">Variables</div>

                <q-input
                  v-model="samplePrinter"
                  dense
                  label="Printer"
                ></q-input>
                <q-input
                  v-model="samplePrintUrl"
                  dense
                  label="URL to print"
                ></q-input>
                <q-input
                  v-model="samplePdfPath"
                  dense
                  label="File path to print"
                ></q-input>
              </q-card-section>
            </q-card>
            <code-block
              v-for="c in s.codeBlocks"
              :label="c.label"
              :language="s.language"
              :code="c.template(snippetTemplateContext)"
              class="relative-position"
            />
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </q-page>
</template>
