import { createRouter, createWebHistory } from "vue-router";
import Home from "@/view/Home.vue";
import Main from "@/view/Main.vue";

const routes = [
  { path: "/", redirect: "/main" },
  {
    path: "/",
    name: "home",
    component: Home,
    children: [{ path: "main", name: "main", component: Main }],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});
export default router;
