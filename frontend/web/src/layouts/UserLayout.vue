<template>
  <div class="user-layout">
    <!-- Header (same style as Home page) -->
    <header class="header" :class="{ 'header-scrolled': scrolled }">
      <div class="header-inner">
        <div class="logo" @click="$router.push('/')">
          <div class="logo-icon">
            <n-icon size="24" color="#fff">
              <AnchorOutline />
            </n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </div>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <router-link to="/user/dashboard" class="nav-link active">控制台</router-link>
        </nav>
        <div class="header-actions">
          <n-badge :value="3" :max="99" :offset="[-4, 0]">
            <n-button text circle size="small" class="notification-btn">
              <template #icon>
                <n-icon :size="18"><NotificationsOutline /></n-icon>
              </template>
            </n-button>
          </n-badge>
          <n-dropdown :options="userDropdownOptions" @select="handleUserAction">
            <div class="user-dropdown-trigger">
              <n-avatar :size="32" class="header-avatar">{{ userInitial }}</n-avatar>
              <span class="header-username">{{ username }}</span>
              <n-icon :size="14" :component="ChevronDownOutline" />
            </div>
          </n-dropdown>
        </div>
      </div>
    </header>

    <!-- Main area: sidebar + content -->
    <div class="layout-body">
      <!-- Mobile overlay -->
      <div v-if="mobileOpen" class="sidebar-overlay" @click="mobileOpen = false"></div>

      <!-- Sidebar -->
      <aside class="sidebar" :class="{ 'mobile-open': mobileOpen }">
        <!-- User card -->
        <div class="user-card">
          <n-avatar :size="64" round class="user-avatar">
            {{ userInitial }}
          </n-avatar>
          <span class="username">{{ username }}</span>
          <n-tag type="info" size="small" round class="user-level">普通用户</n-tag>
        </div>

        <!-- Menu -->
        <nav class="sidebar-menu">
          <router-link
            v-for="item in menuItems"
            :key="item.key"
            :to="item.path"
            class="menu-item"
            :class="{ active: activeMenu === item.key }"
            @click="mobileOpen = false"
          >
            <n-icon :size="20" :component="item.icon" />
            <span class="menu-label">{{ item.label }}</span>
          </router-link>
        </nav>
      </aside>

      <!-- Content -->
      <main class="main-content">
        <div class="page-content">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { DropdownOption } from 'naive-ui'
import { NIcon } from 'naive-ui'
import {
  AnchorOutline,
  HomeOutline,
  CubeOutline,
  ListOutline,
  WalletOutline,
  ChatbubblesOutline,
  TicketOutline,
  PeopleOutline,
  CardOutline,
  PersonOutline,
  ShieldCheckmarkOutline,
  NotificationsOutline,
  ChevronDownOutline,
  PersonCircleOutline,
  SettingsOutline,
  LogOutOutline
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const mobileOpen = ref(false)
const scrolled = ref(false)

const username = computed(() => userStore.username || '用户')
const userInitial = computed(() => username.value.charAt(0).toUpperCase())

function handleScroll() {
  scrolled.value = window.scrollY > 20
}

onMounted(() => window.addEventListener('scroll', handleScroll))
onUnmounted(() => window.removeEventListener('scroll', handleScroll))

const menuItems = [
  { key: 'dashboard', label: '控制台', path: '/user/dashboard', icon: HomeOutline },
  { key: 'products', label: '我的产品', path: '/user/products', icon: CubeOutline },
  { key: 'orders', label: '我的订单', path: '/user/orders', icon: ListOutline },
  { key: 'invoices', label: '我的账单', path: '/user/invoices', icon: WalletOutline },
  { key: 'tickets', label: '工单支持', path: '/user/tickets', icon: ChatbubblesOutline },
  { key: 'coupons', label: '代金券', path: '/user/coupons', icon: TicketOutline },
  { key: 'referral', label: '推荐返利', path: '/user/referral', icon: PeopleOutline },
  { key: 'verification', label: '实名认证', path: '/user/verification', icon: CardOutline },
  { key: 'profile', label: '个人资料', path: '/user/profile', icon: PersonOutline },
  { key: 'security', label: '安全设置', path: '/user/security', icon: ShieldCheckmarkOutline }
]

const routeKeyMap: Record<string, string> = {
  '/user/dashboard': 'dashboard',
  '/user/products': 'products',
  '/user/orders': 'orders',
  '/user/invoices': 'invoices',
  '/user/tickets': 'tickets',
  '/user/coupons': 'coupons',
  '/user/referral': 'referral',
  '/user/verification': 'verification',
  '/user/profile': 'profile',
  '/user/security': 'security'
}

const activeMenu = computed(() => {
  return routeKeyMap[route.path] || 'dashboard'
})

const userDropdownOptions: DropdownOption[] = [
  {
    label: '个人资料',
    key: 'profile',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(PersonCircleOutline) })
  },
  {
    label: '安全设置',
    key: 'security',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(SettingsOutline) })
  },
  { type: 'divider' },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(LogOutOutline) })
  }
]

function handleUserAction(key: string) {
  if (key === 'logout') {
    userStore.logout()
    router.push('/login')
  } else {
    router.push(`/user/${key}`)
  }
}
</script>

<style scoped>
.user-layout {
  min-height: 100vh;
  background: #f0f5ff;
}

/* ==================== Header ==================== */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  z-index: 1000;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
}

.header-scrolled {
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 1px 12px rgba(0, 0, 0, 0.08);
}

.header-inner {
  max-width: 100%;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  user-select: none;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s;
}

.logo:hover .logo-icon {
  transform: rotate(-12deg) scale(1.05);
}

.header-scrolled .logo-icon {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.35);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1px;
  transition: color 0.35s;
}

.header-scrolled .logo-text {
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 36px;
}

.nav-link {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.25s;
  position: relative;
  padding: 4px 0;
}

.nav-link:hover,
.nav-link.active {
  color: #fff;
}

.header-scrolled .nav-link {
  color: #4e5969;
}

.header-scrolled .nav-link:hover,
.header-scrolled .nav-link.active {
  color: #1890ff;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 2px;
  background: #fff;
  border-radius: 1px;
}

.header-scrolled .nav-link.active::after {
  background: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.notification-btn {
  color: rgba(255, 255, 255, 0.85) !important;
}

.header-scrolled .notification-btn {
  color: #595959 !important;
}

.user-dropdown-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-dropdown-trigger:hover {
  background: rgba(255, 255, 255, 0.15);
}

.header-scrolled .user-dropdown-trigger:hover {
  background: #f0f5ff;
}

.header-avatar {
  background: rgba(255, 255, 255, 0.25);
  color: #fff;
  font-weight: 600;
  font-size: 12px;
}

.header-scrolled .header-avatar {
  background: linear-gradient(135deg, #1890ff, #096dd9);
}

.header-username {
  font-size: 14px;
  color: #fff;
  font-weight: 500;
}

.header-scrolled .header-username {
  color: #262626;
}

/* ==================== Layout Body ==================== */
.layout-body {
  display: flex;
  padding-top: 64px;
  min-height: 100vh;
}

/* ==================== Sidebar ==================== */
.sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #f0f0f0;
  position: fixed;
  top: 64px;
  left: 0;
  bottom: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;
}

.user-card {
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.user-avatar {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
  font-weight: 700;
  font-size: 22px;
  flex-shrink: 0;
}

.username {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  text-align: center;
}

.user-level {
  font-size: 12px;
}

/* ==================== Sidebar Menu ==================== */
.sidebar-menu {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  color: #595959;
  text-decoration: none;
  transition: all 0.2s ease;
  font-size: 14px;
  position: relative;
  border: none;
  background: none;
  cursor: pointer;
}

.menu-item:hover {
  background: #f0f5ff;
  color: #1890ff;
}

.menu-item.active {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.menu-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  bottom: 8px;
  width: 3px;
  background: #fff;
  border-radius: 0 2px 2px 0;
}

.menu-item.active :deep(.n-icon) {
  color: #fff;
}

.menu-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ==================== Main Content ==================== */
.main-content {
  flex: 1;
  margin-left: 220px;
  min-height: calc(100vh - 64px);
}

.page-content {
  padding: 24px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.sidebar-overlay {
  display: none;
}

/* ==================== Responsive ==================== */
@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .header-username {
    display: none;
  }

  .sidebar {
    transform: translateX(-100%);
    box-shadow: none;
    z-index: 200;
  }

  .sidebar.mobile-open {
    transform: translateX(0);
    box-shadow: 4px 0 20px rgba(0, 0, 0, 0.1);
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    z-index: 199;
  }

  .main-content {
    margin-left: 0;
  }

  .page-content {
    padding: 16px;
  }
}
</style>
