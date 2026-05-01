<template>
  <el-container class="app-container">
    <el-aside width="200px" class="sidebar">
      <div class="logo">Lattice Coding</div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        router
      >
        <el-menu-item index="/provider">
          <span>模型管理</span>
        </el-menu-item>
        <el-menu-item index="/agent">
          <span>Agent 管理</span>
        </el-menu-item>
        <el-menu-item index="/chat">
          <span>对话</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-main class="content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  background-color: #263445;
}

.sidebar-menu {
  border-right: none;
  background-color: #304156;
}

.sidebar-menu .el-menu-item {
  color: #bfcbd9;
}

.sidebar-menu .el-menu-item:hover,
.sidebar-menu .el-menu-item.is-active {
  background-color: #263445;
  color: #409eff;
}

.content {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
