<template>
  <div class="header-wrap">
    <div class="header">
      <div class="logo">
        <a href="/">{{ site.name }}</a>
      </div>
      <div class="scbox">
        <el-input
          v-model="searchTxt"
          style="height: 30px"
          maxlength="30"
          placeholder="站内搜索"
          value="站内搜索"
          @keyup.enter="searchText"
        >
        </el-input>
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

  <div class="main-wrap">
    <router-view />
    <div class="main" v-if="false">
      <div class="card">
        {{ `template "content" .` }}
      </div>

      <div class="main-sider">
        {{ `template "side" .` }}
      </div>
    </div>
  </div>
  <div class="footer-wrap">
    <div class="footer">
      <p>
        &copy; Copyright <a href="/">{{ site.name }}</a> •
        <a rel="nofollow" href="/feed">Atom Feed</a> •
        <a href="/view?tpl=mobile">手机模式</a>
      </p>
      <p>
        Powered by
        <a href="https://www.youbbs.org" target="_blank">youBBS </a> -
        {{ site.goVersion }}
      </p>
      <p>MD5SUMS: {{ site.MD5Sums }}</p>
    </div>
    <a style="display: none" rel="nofollow" href="#top" id="go-to-top">▲</a>
    <!-- footer end -->
  </div>
</template>

<script>
export default {
  name: "Home",
  data() {
    return {
      site: { goVersion: "1.1", MD5Sums: "" },
      searchTxt: "",
      CurrentUser: {
        id: "",
        name: "",
        noticeNum: 0,
        avatar: "",
        flag: 0,
      },
    };
  },
  methods: {
    searchText() {
      console.log(this.searchTxt);
    },
  },
};
</script>

<style lang="less" scoped>
.header-wrap {
  width: 100%;
  min-width: 1020px;

  color: #fff;
  background: #000;
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
.main-wrap {
  flex: 1;

  width: 100%;
  min-width: 1020px;
  background-color: #ccc;
  padding: 20px 0 20px 0;
}
.hdlogo {
  display: none;
}
.footer-wrap {
  width: 100%;
  padding-bottom: 10px;
  margin-bottom: 0;
  color: #999;
  background-color: #ededed;
  border-top: 1px solid #e0e0e0;
  .footer {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    margin: 0 auto;
    padding: 20px 0 20px 0;
    p {
      margin-bottom: 1px;
      line-height: 120%;
      a {
        color: #606060;
        font-weight: 500;
        &:hover {
          color: #303030;
          text-decoration: underline;
        }
      }
    }
  }
}
</style>
