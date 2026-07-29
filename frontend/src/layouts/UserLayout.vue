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
        <el-avatar :size="48" class="user-avatar">{{ userInitial }}</el-avatar>
        <div class="user-card-info">
          <span class="user-card-name">{{ username }}</span>
          <el-tag type="info" size="small" effect="dark" round>普通用户</el-tag>
        </div>
      </div>

      <!-- Menu -->
      <el-scrollbar class="sidebar-scroll">
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          :collapse="false"
          @select="handleMenuSelect"
        >
          <!-- 快捷入口 -->
          <el-menu-item index="/user/dashboard">
            <el-icon><HomeFilled /></el-icon>
            <span>控制台</span>
          </el-menu-item>
          
          <el-menu-item index="/products">
            <el-icon><ShoppingCart /></el-icon>
            <span>订购产品</span>
          </el-menu-item>
          
          <el-menu-item index="/cart">
            <el-icon><ShoppingCart /></el-icon>
            <span>购物车</span>
          </el-menu-item>
          
          <!-- 业务管理 -->
          <el-sub-menu index="business">
            <template #title>
              <el-icon><Box /></el-icon>
              <span>业务管理</span>
            </template>
            <el-menu-item index="/user/products">我的服务</el-menu-item>
            <el-menu-item index="/user/orders">订单管理</el-menu-item>
            <el-menu-item index="/user/upgrade">产品升降级</el-menu-item>
          </el-sub-menu>
          
          <!-- 财务中心 -->
          <el-sub-menu index="finance">
            <template #title>
              <el-icon><Wallet /></el-icon>
              <span>财务中心</span>
            </template>
            <el-menu-item index="/user/invoices">账单管理</el-menu-item>
            <el-menu-item index="/user/wallet">充值余额</el-menu-item>
            <el-menu-item index="/user/coupons">优惠券</el-menu-item>
          </el-sub-menu>
          
          <!-- 支持服务 -->
          <el-sub-menu index="support">
            <template #title>
              <el-icon><Tickets /></el-icon>
              <span>支持服务</span>
            </template>
            <el-menu-item index="/user/tickets">工单列表</el-menu-item>
            <el-menu-item index="/user/tickets/create">提交工单</el-menu-item>
          </el-sub-menu>
          
          <!-- 资源中心 -->
          <el-sub-menu index="resources">
            <template #title>
              <el-icon><Folder /></el-icon>
              <span>资源中心</span>
            </template>
            <el-menu-item index="/knowledge-base">知识库</el-menu-item>
            <el-menu-item index="/downloads">下载中心</el-menu-item>
            <el-menu-item index="/news">新闻动态</el-menu-item>
          </el-sub-menu>
          
          <!-- 推介计划 -->
          <el-menu-item index="/user/referral">
            <el-icon><Connection /></el-icon>
            <span>推介计划</span>
          </el-menu-item>
          
          <!-- 账户设置 -->
          <el-sub-menu index="account">
            <template #title>
              <el-icon><UserFilled /></el-icon>
              <span>账户设置</span>
            </template>
            <el-menu-item index="/user/profile">个人资料</el-menu-item>
            <el-menu-item index="/user/security">安全设置</el-menu-item>
            <el-menu-item index="/user/verification">实名认证</el-menu-item>
            <el-menu-item index="/user/contacts">联系人管理</el-menu-item>
            <el-menu-item index="/user/oauth-bind">第三方登录</el-menu-item>
          </el-sub-menu>
          
          <!-- 消息中心 -->
          <el-menu-item index="/user/system-message">
            <el-icon><Bell /></el-icon>
            <span>消息中心</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="menu-badge" />
          </el-menu-item>
        </el-menu>
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
            <el-breadcrumb-item>{{ currentPageName }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="top-bar-right">
          <el-badge :value="unreadCount" :max="99" class="notify-badge" :hidden="unreadCount === 0">
            <el-icon :size="18" class="top-icon" @click="$router.push('/user/system-message')"><Bell /></el-icon>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Fold, Close, Bell, ArrowDown, User, Setting, SwitchButton,
  HomeFilled, Box, ShoppingCart, Wallet, Tickets, Ticket, Connection,
  Postcard, UserFilled, Lock, Folder, Download, Document, Promotion, TrendCharts
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const sidebarVisible = ref(false)
const unreadCount = ref(0)

const username = computed(() => userStore.username || '用户')
const userInitial = computed(() => username.value.charAt(0).toUpperCase())

// 当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 当前页面名称
const pageNameMap: Record<string, string> = {
  '/user/dashboard': '控制台',
  '/user/products': '我的服务',
  '/user/orders': '订单管理',
  '/user/upgrade': '产品升降级',
  '/user/invoices': '账单管理',
  '/user/wallet': '充值余额',
  '/user/coupons': '优惠券',
  '/user/tickets': '工单列表',
  '/user/tickets/create': '提交工单',
  '/user/referral': '推介计划',
  '/user/profile': '个人资料',
  '/user/security': '安全设置',
  '/user/verification': '实名认证',
  '/user/contacts': '联系人管理',
  '/user/oauth-bind': '第三方登录',
  '/user/system-message': '消息中心',
  '/user/record-log': '操作日志'
}

const currentPageName = computed(() => {
  return pageNameMap[route.path] || '用户中心'
})

// 菜单选择
const handleMenuSelect = (index: string) => {
  router.push(index)
  sidebarVisible.value = false
}

// 用户操作
const handleUserAction = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else {
    router.push(`/user/${command}`)
  }
}

// 获取未读消息数
const fetchUnreadCount = async () => {
  try {
    const { data } = await request.get('/api/v1/user/messages/unread-count')
    if (data?.data) {
      unreadCount.value = data.data.count || 0
    }
  } catch (error) {
    // 忽略错误
  }
}

onMounted(() => {
  fetchUnreadCount()
})
</script>

<style scoped>
.user-layout {
  display: flex;
  min-height: 100vh;
  background: #f5f7fa;
}

/* ==================== Sidebar ==================== */
.sidebar {
  width: 220px;
  background: #fff;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-right: 1px solid #e8ecf1;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.sidebar-close-mobile {
  display: none;
  color: #909399;
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
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  color: #1a2332;
}

.user-card {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.user-avatar {
  background: linear-gradient(135deg, #409eff, #66b1ff);
  color: #fff;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}

.user-card-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.user-card-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a2332;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Menu */
.sidebar-scroll {
  flex: 1;
  overflow: hidden;
}

.sidebar-menu {
  border-right: none;
}

.sidebar-menu :deep(.el-menu-item) {
  height: 44px;
  line-height: 44px;
  font-size: 14px;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: #ecf5ff;
  color: #409eff;
}

.sidebar-menu :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  font-size: 14px;
}

.menu-badge {
  margin-left: 8px;
}

/* ==================== Main Area ==================== */
.main-area {
  flex: 1;
  margin-left: 220px;
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
  color: #409eff;
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
  background: linear-gradient(135deg, #409eff, #66b1ff);
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
