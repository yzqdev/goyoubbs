import ElementPlus from "element-plus";
import "element-plus/lib/theme-chalk/index.css";
import { createApp } from "vue";
import 'bootstrap/dist/css/bootstrap.min.css'
 import router from "./router";
import App from "./App.vue";
import store from "@/store";


// 设置语言


//全局组件



const app = createApp(App);

//指令

//全局组件
app.use(ElementPlus,{size:'mini',zIndex: 3000 });
app.use(router);
app.use(store);
app.mount("#app");
