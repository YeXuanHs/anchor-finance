<template>
  <el-container class="layout">
    <!-- 顶部导航 -->
    <el-header class="header">
      <div class="header-content">
        <!-- Logo -->
        <div class="logo" @click="router.push('/')">
          <img src="@/assets/logo.png" alt="锚点财务" class="logo-img" />
          <span class="logo-text">{{ siteName }}</span>
        </div>

        <!-- 导航菜单 -->
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          :router="true"
          class="nav-menu"
        >
          <el-menu-item index="/home">首页</el-menu-item>
          <el-menu-item index="/products">产品中心</el-menu-item>
          <el-menu-item index="/services">我的服务</el-menu-item>
          <el-menu-item index="/orders">我的订单</el-menu-item>
          <el-menu-item index="/tickets">工单</el-menu-item>
          <el-menu-item index="/finance">财务中心</el-menu-item>
        </el-menu>

        <!-- 用户信息 -->
        <div class="user-info">
          <template v-if="userStore.isLoggedIn">
            <el-dropdown @command="handleCommand">
              <span class="user-dropdown">
                <el-avatar :size="32" :src="userStore.avatar">
                  {{ userStore.username?.charAt(0) }}
                </el-avatar>
                <span class="username">{{ userStore.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="/user">个人中心</el-dropdown-item>
                  <el-dropdown-item command="/user/profile">个人资料</el-dropdown-item>
                  <el-dropdown-item command="/user/security">安全设置</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="router.push('/login')">登录</el-button>
            <el-button @click="router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 内容区 -->
    <el-main class="main">
      <router-view />
    </el-main>

    <!-- 底部 -->
    <el-footer class="footer">
      <div class="footer-content">
        <p>© {{ currentYear }} {{ siteName }} All Rights Reserved</p>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const currentYear = new Date().getFullYear()
const siteName = import.meta.env.VITE_SITE_NAME || '锚点财务'

const activeMenu = computed(() => {
  return route.path
})

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else {
    router.push(command)
  }
}
</script>

<style scoped lang="scss">
.layout {
  min-height: 100vh;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  height: 60px;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-right: 40px;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.nav-menu {
  flex: 1;
  border-bottom: none;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
  color: #303133;
}

.main {
  background: #f5f7fa;
  min-height: calc(100vh - 120px);
}

.footer {
  background: #545c64;
  color: #fff;
  text-align: center;
  height: 60px;
  line-height: 60px;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;

  p {
    margin: 0;
    font-size: 14px;
  }
}
</style>
