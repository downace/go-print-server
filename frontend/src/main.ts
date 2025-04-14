import { createPinia } from "pinia";
import { Quasar } from "quasar";
import quasarIconSet from "quasar/icon-set/svg-mdi-v7";
import "quasar/dist/quasar.css";
import { createApp } from "vue";
import App from "./App.vue";
import "@quasar/extras/roboto-font/roboto-font.css";
import "@quasar/extras/mdi-v7/mdi-v7.css";
import { router } from "./router";

const app = createApp(App);

app.use(router);
app.use(createPinia());
app.use(Quasar, {
  iconSet: quasarIconSet,
});

app.mount("#app");
