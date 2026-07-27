<template>
  <div class="admin-layout">
    <!-- Sidebar -->
    <aside class="admin-sidebar" :class="{ collapsed: adminStore.sidebarCollapsed }">
      <!-- Logo -->
      <div class="sidebar-logo">
        <el-icon :size="28" color="#409eff">
          <Anchor />
        </el-icon>
        <span v-if="!adminStore.sidebarCollapsed" class="logo-text">AnchorFinance</span>
      </div>

      <!-- Menu -->
      <el-scrollbar class="sidebar-scroll">
        <el-menu
          :default-active="activeMenu"
          :collapse="adminStore.sidebarCollapsed"
          :collapse-transition="false"
          background-color="#001529"
          text-color="rgba(255, 255, 255, 0.65)"
          active-text-color="#ffffff"
          router
        >
          <!-- Dashboard -->
          <el-menu-item index="/admin/dashboard">
            <el-icon><DataBoard /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>

          <!-- User Management -->
          <el-sub-menu index="user-group">
            <template #title>
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </template>
            <el-menu-item index="/admin/users">用户列表</el-menu-item>
            <el-menu-item index="/admin/agents">代理商管理</el-menu-item>
          </el-sub-menu>

          <!-- Business Management -->
          <el-sub-menu index="business-group">
            <template #title>
              <el-icon><ShoppingBag /></el-icon>
              <span>业务管理</span>
            </template>
            <el-menu-item index="/admin/products">产品管理</el-menu-item>
            <el-menu-item index="/admin/orders">订单管理</el-menu-item>
            <el-menu-item index="/admin/coupons">优惠券管理</el-menu-item>
          </el-sub-menu>

          <!-- Finance -->
          <el-sub-menu index="finance-group">
            <template #title>
              <el-icon><Wallet /></el-icon>
              <span>财务管理</span>
            </template>
            <el-menu-item index="/admin/payments">支付管理</el-menu-item>
            <el-menu-item index="/admin/reports">报表统计</el-menu-item>
          </el-sub-menu>

          <!-- Content -->
          <el-sub-menu index="content-group">
            <template #title>
              <el-icon><Document /></el-icon>
              <span>内容管理</span>
            </template>
            <el-menu-item index="/admin/email-templates">邮件模板</el-menu-item>
            <el-menu-item index="/admin/notifications">通知管理</el-menu-item>
          </el-sub-menu>

          <!-- System -->
          <el-sub-menu index="system-group">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/admin/oauth">第三方登录</el-menu-item>
            <el-menu-item index="/admin/logs">操作日志</el-menu-item>
            <el-menu-item index="/admin/settings">系统设置</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </aside>

    <!-- Main Content -->
    <div class="admin-main" :class="{ expanded: adminStore.sidebarCollapsed }">
      <!-- Header -->
      <header class="admin-header">
        <div class="header-left">
          <el-icon
            class="collapse-btn"
            :size="20"
            @click="adminStore.toggleSidebar()"
          >
            <Fold v-if="!adminStore.sidebarCollapsed" />
            <Expand v-else />
          </el-icon>

          <!-- Breadcrumb -->
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentPageTitle">{{ currentPageTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <!-- Notifications -->
          <el-badge :value="3" :max="99" class="notification-badge">
            <el-button :icon="Bell" circle size="small" />
          </el-badge>

          <!-- Fullscreen -->
          <el-button :icon="FullScreen" circle size="small" @click="toggleFullscreen" />

          <!-- User Dropdown -->
          <el-dropdown trigger="click" @command="handleUserCommand">
            <div class="user-info">
              <el-avatar :size="32" :style="{ backgroundColor: '#409eff' }">
                {{ adminStore.adminInfo?.username?.charAt(0)?.toUpperCase() || 'A' }}
              </el-avatar>
              <span class="username">{{ adminStore.adminInfo?.username || '管理员' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :icon="User" command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item :icon="Setting" command="settings">系统设置</el-dropdown-item>
                <el-dropdown-item :icon="SwitchButton" command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Content -->
      <main class="admin-content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Anchor,
  DataBoard,
  User,
  ShoppingBag,
  Wallet,
  Document,
  Setting,
  Fold,
  Expand,
  Bell,
  FullScreen,
  ArrowDown,
  SwitchButton,
} from '@element-plus/icons-vue'
import { useAdminStore } from '~/composables/useAdmin'

const route = useRoute()
const adminStore = useAdminStore()

const activeMenu = computed(() => route.path)

const pageTitleMap: Record<string, string> = {
  '/admin/dashboard': '仪表盘',
  '/admin/users': '用户管理',
  '/admin/agents': '代理商管理',
  '/admin/products': '产品管理',
  '/admin/orders': '订单管理',
  '/admin/coupons': '优惠券管理',
  '/admin/payments': '支付管理',
  '/admin/reports': '报表统计',
  '/admin/email-templates': '邮件模板',
  '/admin/notifications': '通知管理',
  '/admin/oauth': '第三方登录',
  '/admin/logs': '操作日志',
  '/admin/settings': '系统设置',
}

const currentPageTitle = computed(() => pageTitleMap[route.path] || '')

function handleUserCommand(command: string) {
  switch (command) {
    case 'logout':
      adminStore.logout()
      break
    case 'settings':
      navigateTo('/admin/settings')
      break
    case 'profile':
      // Navigate to profile
      break
  }
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

/* Sidebar */
.admin-sidebar {
  width: 240px;
  height: 100vh;
  background: #001529;
  transition: width 0.3s ease;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.admin-sidebar.collapsed {
  width: 64px;
}

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
}

.sidebar-scroll {
  flex: 1;
  overflow: hidden;
}

/* Main */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: margin-left 0.3s ease;
}

/* Header */
.admin-header {
  height: 64px;
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  cursor: pointer;
  color: #666;
  transition: color 0.3s;
}

.collapse-btn:hover {
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.notification-badge {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  font-size: 14px;
  color: #333;
}

/* Content */
.admin-content {
  flex: 1;
  padding: 20px;
  background: #f0f2f5;
  overflow-y: auto;
}
</style>
