<template>
  <div class="user-layout">
    <!-- Mobile Overlay -->
    <div v-if="sidebarVisible" class="sidebar-overlay" @click="sidebarVisible = false" />

    <!-- Sidebar -->
    <aside class="sidebar" :class="{ 'mobile-open': sidebarVisible }">
      <div class="sidebar-header">
        <div class="logo" @click="$router.push('/')">
          <img src="/logo.png" alt="锚点财务" class="logo-img" />
          <span class="logo-text">锚点财务</span>
        </div>
        <el-icon class="sidebar-close-mobile" :size="18" @click="sidebarVisible = false"><Close /></el-icon>
      </div>

      <!-- User Card -->
      <div class="user-card">
        <el-avatar :size="56" class="user-avatar">{{ userInitial }}</el-avatar>
        <div class="user-card-info">
          <span class="user-card-name">{{ username }}</span>
          <el-tag type="info" size="small" effect="dark" round>普通用户</el-tag>
        </div>
      </div>

      <!-- Accordion Menu -->
      <el-scrollbar class="sidebar-scroll">
        <el-collapse v-model="activeMenus" class="menu-collapse">
          <el-collapse-item
            v-for="group in menuGroups"
            :key="group.name"
            :name="group.name"
          >
            <template #title>
              <div class="group-title">
                <el-icon :size="15"><component :is="group.icon" /></el-icon>
                <span>{{ group.label }}</span>
              </div>
            </template>
            <router-link
              v-for="item in group.children"
              :key="item.path"
              :to="item.path"
              class="menu-item"
              :class="{ active: route.path === item.path }"
              @click="sidebarVisible = false"
            >
              <el-icon :size="15"><component :is="item.icon" /></el-icon>
              <span>{{ item.label }}</span>
            </router-link>
          </el-collapse-item>
        </el-collapse>
      </el-scrollbar>
    </aside>

    <!-- Main Area -->
    <div class="main-area">
      <!-- Top Bar -->
      <header class="top-bar">
        <div class="top-bar-left">
          <el-icon class="mobile-menu-btn" :size="20" @click="sidebarVisible = true">
            <Fold />
          </el-icon>
          <!-- Breadcrumb -->
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/user/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentGroup">{{ currentGroup }}</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentPageName }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="top-bar-right">
          <el-badge :value="3" :max="99" class="notify-badge">
            <el-icon :size="18" class="top-icon"><Bell /></el-icon>
          </el-badge>
          <el-dropdown trigger="click" @command="handleUserAction">
            <div class="user-trigger">
              <el-avatar :size="30" class="trigger-avatar">{{ userInitial }}</el-avatar>
              <span class="trigger-name">{{ username }}</span>
              <el-icon :size="12"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人资料
                </el-dropdown-item>
                <el-dropdown-item command="security">
                  <el-icon><Setting /></el-icon>安全设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Content -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Fold, Close, Bell, ArrowDown, User, Setting, SwitchButton,
  HomeFilled, Box, ShoppingCart, Wallet, Tickets, Ticket, Connection,
  Postcard, UserFilled, Lock, Folder, Download, Document, Promotion, TrendCharts
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const sidebarVisible = ref(false)

const username = computed(() => userStore.username || '用户')
const userInitial = computed(() => username.value.charAt(0).toUpperCase())

interface MenuItem {
  path: string
  label: string
  icon: any
}
interface MenuGroup {
  name: string
  label: string
  icon: any
  children: MenuItem[]
}

const menuGroups: MenuGroup[] = [
  {
    name: 'quick',
    label: '快捷入口',
    icon: Promotion,
    children: [
      { path: '/user/dashboard', label: '首页概览', icon: HomeFilled },
      { path: '/products', label: '订购产品', icon: ShoppingCart },
      { path: '/cart', label: '购物车', icon: ShoppingCart },
      { path: '/user/referral', label: '代理升级', icon: TrendCharts }
    ]
  },
  {
    name: 'business',
    label: '业务管理',
    icon: Box,
    children: [
      { path: '/user/products', label: '我的产品', icon: Box },
      { path: '/user/host', label: '我的主机', icon: Box },
      { path: '/user/orders', label: '订单管理', icon: Document },
      { path: '/user/upgrade', label: '升降级', icon: TrendCharts },
      { path: '/user/dcim', label: 'DCIM管理', icon: Box }
    ]
  },
  {
    name: 'finance',
    label: '财务中心',
    icon: Wallet,
    children: [
      { path: '/user/invoices', label: '账单管理', icon: Wallet },
      { path: '/user/wallet', label: '充值余额', icon: Wallet },
      { path: '/user/credit-limit', label: '信用额度', icon: Wallet },
      { path: '/user/coupons', label: '优惠券', icon: Ticket },
      { path: '/user/contract', label: '合同管理', icon: Document }
    ]
  },
  {
    name: 'support',
    label: '支持服务',
    icon: Tickets,
    children: [
      { path: '/user/tickets/create', label: '提交工单', icon: Ticket },
      { path: '/user/tickets', label: '工单列表', icon: Tickets },
      { path: '/knowledge-base', label: '知识库', icon: Folder },
      { path: '/downloads', label: '下载中心', icon: Download }
    ]
  },
  {
    name: 'account',
    label: '账户设置',
    icon: UserFilled,
    children: [
      { path: '/user/profile', label: '个人资料', icon: UserFilled },
      { path: '/user/security', label: '安全设置', icon: Lock },
      { path: '/user/verification', label: '实名认证', icon: Postcard },
      { path: '/user/contacts', label: '联系人管理', icon: UserFilled },
      { path: '/user/oauth-bind', label: '第三方绑定', icon: Connection }
    ]
  },
  {
    name: 'message',
    label: '消息中心',
    icon: Bell,
    children: [
      { path: '/user/system-message', label: '系统消息', icon: Bell },
      { path: '/user/record-log', label: '操作日志', icon: Document }
    ]
  },
  {
    name: 'referral',
    label: '推介计划',
    icon: Connection,
    children: [
      { path: '/user/referral', label: '推介概览', icon: Connection }
    ]
  }
]

function findActiveGroup(): string[] {
  for (const group of menuGroups) {
    if (group.children.some(c => route.path.startsWith(c.path))) {
      return [group.name]
    }
  }
  return ['quick']
}
const activeMenus = ref<string[]>(findActiveGroup())

watch(() => route.path, () => {
  activeMenus.value = findActiveGroup()
})

const currentGroup = computed(() => {
  for (const group of menuGroups) {
    if (group.children.some(c => route.path.startsWith(c.path))) {
      return group.label
    }
  }
  return ''
})

const currentPageName = computed(() => {
  for (const group of menuGroups) {
    const item = group.children.find(c => route.path.startsWith(c.path))
    if (item) return item.label
  }
  return '用户中心'
})

function handleUserAction(command: string) {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else {
    router.push(`/user/${command}`)
  }
}
</script>

<style scoped>
.user-layout {
  display: flex;
  min-height: 100vh;
  background: #f5f7fa;
}

/* ==================== Sidebar ==================== */
.sidebar {
  width: 240px;
  background: #1a3a5c;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255,255,255,0.08);
  flex-shrink: 0;
}

.sidebar-close-mobile {
  display: none;
  color: rgba(255,255,255,0.6);
  cursor: pointer;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}

.logo-img {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: contain;
  transition: transform 0.2s;
}

.logo:hover .logo-img {
  transform: scale(1.08);
}

.logo-text {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.5px;
}

.user-card {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  flex-shrink: 0;
}

.user-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-weight: 700;
  font-size: 18px;
  flex-shrink: 0;
}

.user-card-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.user-card-name {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Accordion Menu */
.sidebar-scroll {
  flex: 1;
  overflow: hidden;
}

.menu-collapse {
  border: none;
}

.menu-collapse :deep(.el-collapse-item__header) {
  height: 42px;
  padding: 0 20px;
  font-size: 12px;
  background: transparent;
  border-bottom: 1px solid rgba(255,255,255,0.04);
  color: rgba(255,255,255,0.45);
}

.menu-collapse :deep(.el-collapse-item__header:hover) {
  color: rgba(255,255,255,0.65);
}

.menu-collapse :deep(.el-collapse-item__wrap) {
  background: transparent;
  border-bottom: none;
}

.menu-collapse :deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

.group-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px 10px 32px;
  color: rgba(255,255,255,0.65);
  text-decoration: none;
  font-size: 13px;
  transition: all 0.2s ease;
  position: relative;
}

.menu-item:hover {
  background: rgba(255,255,255,0.06);
  color: #fff;
}

.menu-item.active {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  color: #fff;
  font-weight: 500;
  border-radius: 0;
}

.menu-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 4px;
  bottom: 4px;
  width: 3px;
  background: #fff;
  border-radius: 0 2px 2px 0;
}

/* ==================== Main Area ==================== */
.main-area {
  flex: 1;
  margin-left: 240px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ==================== Top Bar ==================== */
.top-bar {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e8ecf1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  flex-shrink: 0;
}

.top-bar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.mobile-menu-btn {
  display: none;
  color: #606266;
  cursor: pointer;
}

.top-bar-right {
  display: flex;
  align-items: center;
  gap: 18px;
}

.top-icon {
  color: #606266;
  cursor: pointer;
  transition: color 0.2s;
}

.top-icon:hover {
  color: #0056FF;
}

.notify-badge :deep(.el-badge__content) {
  border: none;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.user-trigger:hover {
  background: #f5f7fa;
}

.trigger-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-weight: 600;
  font-size: 12px;
}

.trigger-name {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.user-trigger .el-icon {
  color: #909399;
}

/* ==================== Main Content ==================== */
.main-content {
  flex: 1;
  padding: 24px;
}

/* ==================== Transitions ==================== */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

.sidebar-overlay {
  display: none;
}

/* ==================== Responsive ==================== */
@media (max-width: 768px) {
  .mobile-menu-btn {
    display: flex;
  }

  .sidebar {
    transform: translateX(-100%);
    z-index: 200;
    transition: transform 0.3s ease;
  }

  .sidebar.mobile-open {
    transform: translateX(0);
    box-shadow: 4px 0 24px rgba(0,0,0,0.2);
  }

  .sidebar-close-mobile {
    display: block;
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.45);
    z-index: 199;
  }

  .main-area {
    margin-left: 0;
  }

  .main-content {
    padding: 16px;
  }

  .top-bar {
    padding: 0 16px;
  }

  .trigger-name {
    display: none;
  }
}
</style>
