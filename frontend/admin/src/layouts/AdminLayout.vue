<template>
  <el-container class="admin-layout">
    <!-- Sidebar -->
    <el-aside :width="sidebarWidth" class="admin-sidebar">
      <div class="sidebar-logo" :class="{ collapsed: isCollapsed }">
        <el-icon :size="28" color="#4080FF">
          <Anchor />
        </el-icon>
        <span v-if="!isCollapsed" class="logo-text">智简魔方</span>
      </div>
      <el-scrollbar>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapsed"
          :collapse-transition="false"
          background-color="#1a3a5c"
          text-color="rgba(255,255,255,0.7)"
          active-text-color="#ffffff"
          router
        >
          <template v-for="group in menuGroups" :key="group.key">
            <el-menu-item-group>
              <template #title>
                <span class="menu-group-title">{{ group.label }}</span>
              </template>
              <el-menu-item
                v-for="item in group.children"
                :key="item.path"
                :index="item.path"
              >
                <el-icon><component :is="item.icon" /></el-icon>
                <span>{{ item.label }}</span>
              </el-menu-item>
            </el-menu-item-group>
          </template>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <!-- Main Area -->
    <el-container class="admin-main-container">
      <!-- Header -->
      <el-header class="admin-header" height="64px">
        <div class="header-left">
          <el-icon
            class="collapse-trigger"
            :size="20"
            @click="toggleSidebar"
          >
            <Fold v-if="!isCollapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/" class="admin-breadcrumb">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-tooltip content="全屏" placement="bottom">
            <el-icon class="header-action" :size="18" @click="toggleFullscreen">
              <FullScreen />
            </el-icon>
          </el-tooltip>
          <el-badge :value="3" :max="99" class="header-badge">
            <el-icon class="header-action" :size="18">
              <Bell />
            </el-icon>
          </el-badge>
          <el-dropdown trigger="click" @command="handleUserAction">
            <div class="header-user">
              <el-avatar :size="32" :style="{ backgroundColor: '#0056FF' }">
                {{ adminStore.adminInfo?.username?.charAt(0)?.toUpperCase() || 'A' }}
              </el-avatar>
              <span class="username">{{ adminStore.adminInfo?.username || '管理员' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- Content -->
      <el-main class="admin-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import {
  DataBoard,
  User,
  Goods,
  ShoppingCart,
  Document,
  ChatDotSquare,
  Ticket,
  Megaphone,
  Setting,
  Key,
  Postcard,
  Operation,
  Bell,
  Fold,
  Expand,
  FullScreen,
  ArrowDown,
  SwitchButton,
  Anchor,
  TrendCharts,
  Avatar,
  Message,
  Notification,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const isCollapsed = computed(() => adminStore.sidebarCollapsed)
const sidebarWidth = computed(() => (isCollapsed.value ? '64px' : '220px'))
const activeMenu = computed(() => route.path)
const currentTitle = computed(() => (route.meta?.title as string) || '')

interface MenuItem {
  label: string
  path: string
  icon: any
}

interface MenuGroup {
  key: string
  label: string
  children: MenuItem[]
}

const menuGroups: MenuGroup[] = [
  {
    key: 'home',
    label: '主页',
    children: [
      { label: '仪表盘', path: '/admin/dashboard', icon: DataBoard },
    ],
  },
  {
    key: 'users',
    label: '用户管理',
    children: [
      { label: '用户管理', path: '/admin/users', icon: User },
      { label: '代理商管理', path: '/admin/agents', icon: Avatar },
    ],
  },
  {
    key: 'business',
    label: '业务管理',
    children: [
      { label: '产品管理', path: '/admin/products', icon: Goods },
      { label: '订单管理', path: '/admin/orders', icon: ShoppingCart },
      { label: '账单管理', path: '/admin/invoices', icon: Document },
      { label: '工单管理', path: '/admin/tickets', icon: ChatDotSquare },
      { label: '优惠券管理', path: '/admin/coupons', icon: Ticket },
    ],
  },
  {
    key: 'content',
    label: '内容管理',
    children: [
      { label: '公告管理', path: '/admin/announcements', icon: Megaphone },
      { label: '邮件模板', path: '/admin/email-templates', icon: Message },
      { label: '通知管理', path: '/admin/notifications', icon: Notification },
    ],
  },
  {
    key: 'system',
    label: '系统管理',
    children: [
      { label: '支付管理', path: '/admin/payments', icon: Postcard },
      { label: '第三方登录', path: '/admin/oauth', icon: Key },
      { label: '报表统计', path: '/admin/reports', icon: TrendCharts },
      { label: '系统日志', path: '/admin/logs', icon: Operation },
      { label: '系统设置', path: '/admin/settings', icon: Setting },
    ],
  },
]

function toggleSidebar() {
  adminStore.toggleSidebar()
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

function handleUserAction(command: string) {
  if (command === 'logout') {
    adminStore.logout()
  }
}

onMounted(() => {
  adminStore.fetchAdminInfo()
})
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  overflow: hidden;
}

/* ── Sidebar ── */
.admin-sidebar {
  background: #1a3a5c;
  transition: width 0.3s ease;
  overflow: hidden;
}

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

.sidebar-logo.collapsed {
  padding: 0;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
  letter-spacing: 1px;
}

.menu-group-title {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.admin-sidebar :deep(.el-menu) {
  border-right: none;
}

.admin-sidebar :deep(.el-menu-item) {
  height: 44px;
  line-height: 44px;
  margin: 2px 8px;
  border-radius: 6px;
}

.admin-sidebar :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%) !important;
  color: #fff !important;
  font-weight: 500;
}

.admin-sidebar :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.06) !important;
}

.admin-sidebar :deep(.el-menu-item-group__title) {
  padding: 12px 0 4px 20px;
}

/* ── Header ── */
.admin-header {
  background: #1a3a5c;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-trigger {
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: color 0.2s;
}

.collapse-trigger:hover {
  color: #fff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-action {
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: color 0.2s;
}

.header-action:hover {
  color: #fff;
}

.header-badge {
  line-height: 1;
}

.header-user {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.85);
}

.username {
  font-size: 14px;
  font-weight: 500;
}

/* ── Main Content ── */
.admin-main-container {
  overflow: hidden;
}

.admin-content {
  background: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}
</style>
