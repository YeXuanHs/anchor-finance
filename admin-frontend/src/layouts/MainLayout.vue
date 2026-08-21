<template>
  <a-layout class="layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      :collapsed="collapsed"
      :width="220"
      :collapsed-width="64"
      breakpoint="xl"
      @collapse="onCollapse"
      class="sider"
    >
      <!-- Logo -->
      <div class="logo">
        <img src="@/assets/logo.png" alt="锚点财务" class="logo-img" />
        <span v-show="!collapsed" class="logo-text">锚点财务</span>
      </div>

      <!-- 菜单 -->
      <a-menu
        :collapsed="collapsed"
        :default-open-keys="['customer', 'order', 'finance', 'ticket', 'feature', 'setting']"
        :selected-keys="selectedKeys"
        @menu-item-click="onMenuClick"
        @sub-menu-click="onSubMenuClick"
      >
        <!-- 仪表盘 -->
        <a-menu-item key="/dashboard">
          <template #icon><icon-dashboard /></template>
          仪表盘
        </a-menu-item>

        <!-- 客户 -->
        <a-sub-menu key="customer">
          <template #icon><icon-user /></template>
          <template #title>客户</template>
          <a-menu-item key="/customer/list">客户列表</a-menu-item>
          <a-menu-item key="/customer/authentication">实名认证</a-menu-item>
          <a-menu-item key="/customer/resources">客户资源池</a-menu-item>
        </a-sub-menu>

        <!-- 业务 -->
        <a-sub-menu key="order">
          <template #icon><icon-file /></template>
          <template #title>业务</template>
          <a-menu-item key="/order/list">产品订单</a-menu-item>
          <a-menu-item key="/order/renewal">续费订单</a-menu-item>
          <a-menu-item key="/order/service">业务列表</a-menu-item>
        </a-sub-menu>

        <!-- 财务 -->
        <a-sub-menu key="finance">
          <template #icon><icon-money-circle /></template>
          <template #title>财务</template>
          <a-menu-item key="/finance/transactions">交易流水</a-menu-item>
          <a-menu-item key="/finance/invoices">账单管理</a-menu-item>
          <a-menu-item key="/finance/credit">信用额管理</a-menu-item>
        </a-sub-menu>

        <!-- 工单 -->
        <a-sub-menu key="ticket">
          <template #icon><icon-customer-service /></template>
          <template #title>工单</template>
          <a-menu-item key="/ticket/list">工单列表</a-menu-item>
          <a-menu-item key="/ticket/statistics">工单统计</a-menu-item>
        </a-sub-menu>

        <!-- 功能 -->
        <a-sub-menu key="feature">
          <template #icon><icon-apps /></template>
          <template #title>功能</template>
          <a-menu-item key="/feature/plugins">插件列表</a-menu-item>
          <a-menu-item key="/feature/statistics">统计</a-menu-item>
        </a-sub-menu>

        <!-- 设置 -->
        <a-sub-menu key="setting">
          <template #icon><icon-settings /></template>
          <template #title>设置</template>
          <a-menu-item key="/setting/general">常规设置</a-menu-item>
          <a-menu-item key="/setting/payment">支付接口</a-menu-item>
          <a-menu-item key="/setting/security">安全相关</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <!-- 内容区 -->
    <a-layout>
      <!-- 顶部导航 -->
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <template #icon>
              <icon-menu-fold v-if="!collapsed" />
              <icon-menu-unfold v-else />
            </template>
          </a-button>
          <!-- 面包屑 -->
          <a-breadcrumb class="breadcrumb">
            <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
              {{ item.title }}
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <!-- 用户信息 -->
          <a-dropdown>
            <a-button type="text">
              <icon-user />
              <span class="username">{{ username }}</span>
            </a-button>
            <template #content>
              <a-doption @click="handleLogout">退出登录</a-doption>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 内容 -->
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const collapsed = ref(false)
const username = computed(() => userStore.username || '管理员')

const selectedKeys = computed(() => {
  return [route.path]
})

const breadcrumbs = computed(() => {
  const matched = route.matched
  return matched.map(item => ({
    path: item.path,
    title: item.meta.title || ''
  }))
})

const onCollapse = (val: boolean) => {
  collapsed.value = val
}

const onMenuClick = (key: string) => {
  router.push(key)
}

const onSubMenuClick = (key: string) => {
  // 子菜单点击
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped lang="scss">
.layout {
  height: 100vh;
}

.sider {
  background: #001529;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-img {
  width: 32px;
  height: 32px;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  margin-left: 12px;
  white-space: nowrap;
}

.header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.breadcrumb {
  margin-left: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.username {
  margin-left: 8px;
}

.content {
  padding: 24px;
  background: #f5f5f5;
  min-height: calc(100vh - 64px);
}
</style>
