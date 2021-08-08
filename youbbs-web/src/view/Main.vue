<template>
 <div class="main-container">

   <div class="main-left">
     <el-card >
       <template #header>
         <div class="card-header">
           <span>卡片名称</span>
           <el-button   type="primary">+发新帖</el-button>
         </div>
       </template>
     </el-card>
   </div>
   <div class="main-right">

     <el-card>
       <template #header>管理员面板</template>
       <div class="add">
         <a class="btn btn-primary" href="/admin/category/list">分类管理</a>
         <a class="btn btn-primary" href="/admin/user/list">用户管理</a>
         <a class="btn btn-primary" href="/admin/link/list">链接管理</a>
       </div>
     </el-card>
   </div>

 </div>
</template>

<script>
export default {
  name: "Main"
}
</script>

<style lang="less" scoped>

  .main-container {
    width: 100%;
    padding: 20px;
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 20px;
    justify-content: center;
    .main-left{
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
    }
    .main-right{
      .main-sider {
        .card:not(:first-of-type) {
          margin-top: 20px;
        }
      }
    }
  }


.main-content {
  width: 730px;
}

</style>
