<template>
  <el-container class="basic-layout">
    <el-aside :width="appStore.collapsed ? '64px' : '220px'" class="sidebar sidebar-dark">
      <div class="logo-area">
        <div class="logo" :class="{ 'logo-collapsed': appStore.collapsed }">
          <div v-if="appStore.collapsed" class="logo-avatar">
            <span>L</span>
          </div>
          <template v-else>
            <span class="logo-text">Lattice-Coding</span>
            <span class="logo-subtitle">AI Agent Platform</span>
          </template>
        </div>
      </div>

      <el-menu
        :default-active="currentRoute"
        :collapse="appStore.collapsed"
        :collapse-transition="false"
        class="sidebar-menu"
        router
      >
        <el-menu-item index="/provider">
          <el-icon><Setting /></el-icon>
          <template #title>模型管理</template>
        </el-menu-item>
        <el-menu-item index="/agent">
          <el-icon><User /></el-icon>
          <template #title>Agent 管理</template>
        </el-menu-item>
        <el-menu-item index="/chat">
          <el-icon><ChatDotRound /></el-icon>
          <template #title>对话</template>
        </el-menu-item>
        <el-menu-item index="/run">
          <el-icon><VideoPlay /></el-icon>
          <template #title>Run</template>
        </el-menu-item>
        <el-menu-item index="/knowledge">
          <el-icon><Collection /></el-icon>
          <template #title>知识库</template>
        </el-menu-item>
        <el-menu-item index="/workflow">
          <el-icon><Connection /></el-icon>
          <template #title>工作流</template>
        </el-menu-item>
        <el-menu-item index="/mcp">
          <el-icon><Grid /></el-icon>
          <template #title>MCP</template>
        </el-menu-item>
        <el-menu-item index="/safety">
          <el-icon><Lock /></el-icon>
          <template #title>安全</template>
        </el-menu-item>
      </el-menu>

      <div class="sidebar-footer">
        <div class="version-info" v-if="!appStore.collapsed">v1.0.0</div>
        <el-button
          text
          class="collapse-btn"
          @click="appStore.toggleCollapse"
        >
          <el-icon v-if="appStore.collapsed"><DArrowRight /></el-icon>
          <el-icon v-else><DArrowLeft /></el-icon>
        </el-button>
      </div>
    </el-aside>

    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
              {{ item.meta.title || item.label }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-icon class="header-icon"><Bell /></el-icon>
          <div class="user-info">
            <el-avatar :size="32" class="user-avatar">{{ getUserInitial('Admin') }}</el-avatar>
            <span class="user-name">Admin</span>
          </div>
        </div>
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
import {
  Setting,
  User,
  ChatDotRound,
  VideoPlay,
  Collection,
  Connection,
  Grid,
  Lock,
  DArrowLeft,
  DArrowRight,
  Bell
} from '@element-plus/icons-vue'

const route = useRoute()
const appStore = useAppStore()

const currentRoute = computed(() => route.path)

const breadcrumbs = computed(() => {
  return route.matched.filter(item => item.meta && item.meta.title)
    .map(item => ({
      path: item.path,
      meta: item.meta as { title?: string },
      label: item.meta.title || ''
    }))
})

const getUserInitial = (name: string): string => {
  if (!name) return 'U'
  return name.charAt(0).toUpperCase()
}
</script>

<style scoped>
.basic-layout {
  height: 100vh;
}

.sidebar {
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  background: var(--color-bg-dark) !important;
  border-right: 1px solid var(--color-border-sidebar);
}

.logo-area {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-border-sidebar);
}

.logo {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.logo-collapsed {
  gap: 0;
}

.logo-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary-500), var(--color-accent-500));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-primary-400), var(--color-accent-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
}

.logo-collapsed .logo-text {
  font-size: 20px;
}

.logo-subtitle {
  font-size: 10px;
  color: var(--color-text-sidebar-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.logo-collapsed .logo-subtitle {
  display: none;
}

.sidebar-menu {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
  background: transparent !important;
  border: none;
}

.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

.sidebar-menu .el-menu-item {
  color: rgba(255, 255, 255, 0.85) !important;
  margin: 4px 8px;
  border-radius: 8px;
  height: 44px;
  transition: all var(--transition-normal);
  background: transparent !important;
}

.sidebar-menu .el-menu-item:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #fff !important;
}

.sidebar-menu .el-menu-item.is-active {
  background: rgba(139, 92, 246, 0.15) !important;
  color: #fff !important;
  position: relative;
}

.sidebar-menu .el-menu-item.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(180deg, var(--color-primary-400), var(--color-primary-600));
  border-radius: 0 2px 2px 0;
}

.sidebar-footer {
  padding: 12px 8px;
  border-top: 1px solid var(--color-border-sidebar);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.version-info {
  font-size: 11px;
  color: var(--color-text-sidebar-muted);
  padding-left: 8px;
}

.collapse-btn {
  color: var(--color-text-sidebar-muted) !important;
  transition: all var(--transition-normal);
}

.collapse-btn:hover {
  color: #fff !important;
  background: rgba(255, 255, 255, 0.1) !important;
}

.main-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--color-bg-secondary);
}

.header {
  background-color: var(--color-bg-primary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  font-size: 18px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all var(--transition-fast);
}

.header-icon:hover {
  background: var(--color-bg-hover);
  color: var(--color-text-primary);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all var(--transition-fast);
}

.user-info:hover {
  background: var(--color-bg-hover);
}

.user-avatar {
  background: linear-gradient(135deg, var(--color-primary-500), var(--color-primary-600));
  color: #fff;
  font-size: 12px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.content {
  background-color: var(--color-bg-secondary);
  padding: 24px;
  overflow-y: auto;
}
</style>
