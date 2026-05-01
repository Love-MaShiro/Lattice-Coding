<template>
  <el-container class="basic-layout">
    <el-aside :width="appStore.collapsed ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <span v-if="!appStore.collapsed">Lattice Coding</span>
        <span v-else>LC</span>
      </div>
      <el-menu
        :default-active="currentRoute"
        :collapse="appStore.collapsed"
        class="sidebar-menu"
        router
      >
        <el-menu-item index="/provider">
          <span>Provider</span>
        </el-menu-item>
        <el-menu-item index="/agent">
          <span>Agent</span>
        </el-menu-item>
        <el-menu-item index="/chat">
          <span>Chat</span>
        </el-menu-item>
        <el-menu-item index="/run">
          <span>Run</span>
        </el-menu-item>
        <el-menu-item index="/knowledge">
          <span>Knowledge</span>
        </el-menu-item>
        <el-menu-item index="/workflow">
          <span>Workflow</span>
        </el-menu-item>
        <el-menu-item index="/mcp">
          <span>MCP</span>
        </el-menu-item>
        <el-menu-item index="/safety">
          <span>Safety</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-title">{{ appStore.title }}</div>
      </el-header>
      <el-main class="content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'

const route = useRoute()
const appStore = useAppStore()

const currentRoute = computed(() => route.path)
</script>

<style scoped>
.basic-layout {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
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

.header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: 500;
}

.content {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
