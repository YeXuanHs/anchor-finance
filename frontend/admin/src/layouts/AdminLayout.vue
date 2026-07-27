<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      bordered
      :collapsed="adminStore.sidebarCollapsed"
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      show-trigger
      @collapse="adminStore.toggleSidebar()"
      @expand="adminStore.toggleSidebar()"
      :native-scrollbar="false"
      style="background: #001529"
    >
      <div class="logo" :class="{ collapsed: adminStore.sidebarCollapsed }">
        <n-icon size="28" color="#18a058">
          <AnchorIcon />
        </n-icon>
        <span v-if="!adminStore.sidebarCollapsed" class="logo-text">锚点财务</span>
      </div>
      <n-menu
        :collapsed="adminStore.sidebarCollapsed"
        :collapsed-width="64"
        :collapsed-icon-size="20"
        :options="menuOptions"
        :value="activeMenu"
        @update:value="handleMenuClick"
        inverted
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered style="height: 64px; display: flex; align-items: center; justify-content: space-between; padding: 0 24px">
        <div style="display: flex; align-items: center; gap: 16px">
          <n-breadcrumb>
            <n-breadcrumb-item>
              <router-link to="/admin/dashboard">首页</router-link>
            </n-breadcrumb-item>
            <n-breadcrumb-item v-if="currentTitle">
              {{ currentTitle }}
            </n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <div style="display: flex; align-items: center; gap: 16px">
          <n-badge :value="3" :max="99">
            <n-button quaternary circle>
              <template #icon>
                <n-icon><NotificationsIcon /></n-icon>
              </template>
            </n-button>
          </n-badge>
          <n-dropdown :options="userDropdownOptions" @select="handleUserAction">
            <div style="display: flex; align-items: center; gap: 8px; cursor: pointer">
              <n-avatar round :size="32" :style="{ backgroundColor: '#18a058' }">
                {{ adminStore.adminInfo?.username?.charAt(0)?.toUpperCase() || 'A' }}
              </n-avatar>
              <span>{{ adminStore.adminInfo?.username || '管理员' }}</span>
            </div>
          </n-dropdown>
        </div>
      </n-layout-header>
      <n-layout-content content-style="padding: 24px;" :native-scrollbar="false">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted, type Component } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { NIcon } from 'naive-ui'
import type { MenuOption, DropdownOption } from 'naive-ui'
import {
  GridOutline as DashboardIcon,
  PeopleOutline as UsersIcon,
  CubeOutline as ProductsIcon,
  CartOutline as OrdersIcon,
  ReceiptOutline as InvoicesIcon,
  ChatbubblesOutline as TicketsIcon,
  MegaphoneOutline as AnnouncementsIcon,
  SettingsOutline as SettingsIcon,
  LogOutOutline as LogoutIcon,
  PersonOutline as PersonIcon,
  NotificationsOutline as NotificationsIcon,
  AccessibilityOutline as AnchorIcon,
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => (route.meta?.title as string) || '')

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  {
    label: '仪表盘',
    key: '/admin/dashboard',
    icon: renderIcon(DashboardIcon),
  },
  {
    label: '用户管理',
    key: '/admin/users',
    icon: renderIcon(UsersIcon),
  },
  {
    label: '产品管理',
    key: '/admin/products',
    icon: renderIcon(ProductsIcon),
  },
  {
    label: '订单管理',
    key: '/admin/orders',
    icon: renderIcon(OrdersIcon),
  },
  {
    label: '账单管理',
    key: '/admin/invoices',
    icon: renderIcon(InvoicesIcon),
  },
  {
    label: '工单管理',
    key: '/admin/tickets',
    icon: renderIcon(TicketsIcon),
  },
  {
    label: '公告管理',
    key: '/admin/announcements',
    icon: renderIcon(AnnouncementsIcon),
  },
  {
    label: '系统设置',
    key: '/admin/settings',
    icon: renderIcon(SettingsIcon),
  },
]

const userDropdownOptions: DropdownOption[] = [
  {
    label: '个人信息',
    key: 'profile',
    icon: renderIcon(PersonIcon),
  },
  {
    type: 'divider',
    key: 'd1',
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogoutIcon),
  },
]

function handleMenuClick(key: string) {
  router.push(key)
}

function handleUserAction(key: string) {
  if (key === 'logout') {
    adminStore.logout()
  }
}

onMounted(() => {
  adminStore.fetchAdminInfo()
})
</script>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo.collapsed {
  padding: 0;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
}
</style>
