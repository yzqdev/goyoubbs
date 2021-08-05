import { createStore } from "vuex";

const store = createStore({
    state() {
        return {
            siteConfig: null,

        };
    },
    mutations: {
        siteConfig(state, siteConfig) {
            state.siteConfig = siteConfig;
        },

    },
});
export default store;
