<template>
  <div class="app">
    <router-view />
  </div>
</template>

<script>
export default {
  data() {
    return {
      visible: false,
    };
  },
  computed: {},
  mounted() {},
  methods: {},
};
</script>
<style lang="less">
html,body{
  margin: 0;
  padding: 0;
}
.app {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
}
</style>
