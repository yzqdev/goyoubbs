<template>
  <div class="header-wrap">
    <div class="header">
      <div class="logo">
        <a href="/">{{ site.name }}</a>
      </div>
      <div class="scbox">
        <form
          role="search"
          method="get"
          id="searchform"
          onsubmit="return dispatch()"
          target="_blank"
        >
          <input
            type="text"
            class="form-control"
            style="height: 30px"
            maxlength="30"
            onfocus="if(this.value=='站内搜索') this.value='';"
            onblur="if(this.value=='') this.value='站内搜索';"
            value="站内搜索"
            name="q"
            id="q"
          />
        </form>
      </div>

      <div class="banner">
        <template v-if="CurrentUser.id">
          <img
            class="avatar avatar24"
            :src="`/static/avatar/${CurrentUser.avatar}.jpg`"
            :alt="CurrentUser.name"
          />
        </template>
        <!--        {{if not .CurrentUser.Password }}-->
        <!--        {{ end }}-->
        <template v-if="CurrentUser.noticeNum > 0">
          <a href="/notification" style="color: yellow"
            >{{ CurrentUser.noticeNum }}条提醒</a
          ></template
        >

        <template v-if="CurrentUser.flag == 0"
          ><span style="color: yellow">已被禁用</span
          >&nbsp;&nbsp;&nbsp;</template
        >
        <template v-if="CurrentUser.flag == 1">
          <span style="color: yellow">在等待审核</span
          >&nbsp;&nbsp;&nbsp;</template
        >

        <a href="/member/{{.CurrentUser.Id}}">{{ CurrentUser.name }}</a
        >&nbsp;&nbsp;&nbsp;<a href="/setting">设置</a>&nbsp;&nbsp;&nbsp;<a
          href="/logout"
          >退出</a
        >

        <template v-if="site.weiboClientId">
          <a href="/wblogin" rel="nofollow"
            ><img
              src="/static/img/weibo_login_55_24.png"
              alt="微博登录"
              title="用微博帐号登录" /></a
          >&nbsp;&nbsp;
        </template>
        <template v-if="site.qqClientId">
          <a href="/qqlogin" rel="nofollow"
            ><img
              src="/static/img/qq_login_55_24.png"
              alt="QQ登录"
              title="用QQ登录" /></a
          >&nbsp;&nbsp;
        </template>

        <template v-if="!site.closeReg"
          >&nbsp;&nbsp;<a href="/register" rel="nofollow">注册</a></template
        >
        <a href="/login" rel="nofollow">登录</a>
      </div>
    </div>
    <!-- header end -->
  </div>
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
  background: #000 url("~@assets/img/bg_header.png") repeat-x bottom center;
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
