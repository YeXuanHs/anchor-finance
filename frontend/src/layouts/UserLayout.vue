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

      <!-- Menu (动态从数据库读取) -->
      <el-scrollbar class="sidebar-scroll">
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          :collapse="false"
          @select="handleMenuSelect"
        >
          <template v-for="item in menuList" :key="item.id">
            <!-- 有子菜单的 -->
            <el-sub-menu v-if="item.children && item.children.length > 0" :index="'menu-' + item.id">
              <template #title>
                <el-icon v-if="item.fa_icon"><component :is="getIconComponent(item.fa_icon)" /></el-icon>
                <span>{{ item.name }}</span>
              </template>
              <el-menu-item 
                v-for="child in item.children" 
                :key="child.id" 
                :index="child.url"
              >
                {{ child.name }}
              </el-menu-item>
            </el-sub-menu>

            <!-- 没有子菜单的 -->
            <el-menu-item v-else :index="item.url">
              <el-icon v-if="item.fa_icon"><component :is="getIconComponent(item.fa_icon)" /></el-icon>
              <span>{{ item.name }}</span>
            </el-menu-item>
          </template>
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
import { ref, computed, onMounted, markRaw } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Fold, Close, Bell, ArrowDown, User, Setting, SwitchButton,
  HomeFilled, Box, ShoppingCart, Wallet, Tickets, Ticket, Connection,
  Postcard, UserFilled, Lock, Folder, Download, Document, Promotion, TrendCharts, Shop
} from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const sidebarVisible = ref(false)
const unreadCount = ref(0)
const menuList = ref<any[]>([])

const username = computed(() => userStore.username || '用户')
const userInitial = computed(() => username.value.charAt(0).toUpperCase())

// 图标映射
const iconMap: Record<string, any> = {
  'bx bx-home-circle': markRaw(HomeFilled),
  'bx bxs-grid-alt': markRaw(Box),
  'bx bx-user': markRaw(UserFilled),
  'bx bx-dollar-circle': markRaw(Wallet),
  'bx bx-detail': markRaw(Tickets),
  'bx bxs-paper-plane': markRaw(Connection),
  'bx bx-store': markRaw(Shop),
  'HomeFilled': markRaw(HomeFilled),
  'ShoppingCart': markRaw(ShoppingCart),
  'Box': markRaw(Box),
  'Wallet': markRaw(Wallet),
  'Tickets': markRaw(Tickets),
  'Folder': markRaw(Folder),
  'Connection': markRaw(Connection),
  'UserFilled': markRaw(UserFilled),
  'Bell': markRaw(Bell),
  'Shop': markRaw(Shop),
}

function getIconComponent(icon: string) {
  return iconMap[icon] || markRaw(Box)
}

// 当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 当前页面名称
const currentPageName = computed(() => {
  const findName = (items: any[]): string => {
    for (const item of items) {
      if (item.url === route.path) return item.name
      if (item.children) {
        const found = findName(item.children)
        if (found) return found
      }
    }
    return ''
  }
  return findName(menuList.value) || '用户中心'
})

// 获取菜单
const fetchMenus = async () => {
  try {
    const { data } = await request.get('/v1/user/menus')
    if (data) {
      menuList.value = Array.isArray(data) ? data : []
    }
  } catch (error) {
    console.error('获取菜单失败:', error)
  }
}

// 菜单选择
const handleMenuSelect = (index: string) => {
  if (index.startsWith('/')) {
    router.push(index)
  }
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
  fetchMenus()
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
