import CodeSnippets from "@/components/CodeSnippets.vue";
import { createRouter, createWebHashHistory } from "vue-router";
import AppDashboard from "./components/AppDashboard.vue";
import AppSettings from "./components/AppSettings.vue";

const router = createRouter({
  // do not use `createWebHistory`: https://github.com/wailsapp/wails/issues/2262
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: AppDashboard,
    },
    {
      path: "/settings",
      component: AppSettings,
    },
    {
      path: "/snippets",
      component: CodeSnippets,
    },
  ],
});

export { router };
