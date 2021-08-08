<template>

</template>

<script>
//   var dispatch = function () {
//   q = document.getElementById("q");
//   if (q.value != "" && q.value != "站内搜索") {
//   window.location.href = "/search?q=" + q.value;
//   return false;
// } else {
//   return false;
// }
// };

export default {
  name: "Header",
  data() {
    return {
      site: {},
      CurrentUser: {
        id: "",
        name: "",
        noticeNum: 0,
        avatar: "",
        flag: 0,
      },
    };
  },
};
</script>

<style lang="less" scoped>
.header-wrap {
  width: 100%;
  min-width: 1020px;

  color: #fff;
  background: #000 url("@/assets/img/bg_header.png") repeat-x bottom center;
  padding: 0;
  -webkit-box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.4);
  .header {
    padding: 0 100px;
    display: flex;
    align-items: center;
    margin: 0 auto;
    color: #666;
    vertical-align: middle;
    font-size: 15px;
    line-height: 19px;
    font-weight: 500;
    .scbox {
      float: left;
      width: 200px;
    }

    .banner {
      flex: 1;
      text-align: right;
      padding: 15px;
      img {
        vertical-align: middle;
        border-radius: 4px;
      }
    }
    .logo {
      display: block;
      float: left;
      font-size: 30px;
      line-height: 40px;

      a {
        text-decoration: none;
        color: #ccc;
        text-shadow: #000 0 0 20px;
        -webkit-transition: color 0.25s linear;
        margin-right: 20px;
        &:hover,
        &:focus,
        &:active {
          color: #fff;
          text-shadow: #666 0 0 20px;
        }
      }
    }
    a {
      color: #ccc;
      text-decoration: none;
      &:hover {
        color: #fff;
        text-decoration: none;
      }
    }
  }
}

.hdlogo {
  display: none;
}
</style>
