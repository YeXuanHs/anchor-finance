<template>
  <header class="site-header" :class="{ 'header-scrolled': scrolled }">
    <div class="header-inner">
      <div class="logo" @click="$router.push('/')">
        <img :src="configStore.getLogo('home') || '/logo.png'" :alt="siteName" class="logo-img" />
        <span class="logo-text">{{ siteName }}</span>
      </div>
      
      <nav class="nav-links">
        <router-link to="/" class="nav-link">首页</router-link>
        
        <!-- 产品下拉菜单 -->
        <el-dropdown trigger="hover" @command="(cmd: string) => $router.push(`/products?group=${cmd}`)">
          <span class="nav-link">
            产品<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="group in productGroups" :key="group.id" :command="group.id">
                {{ group.name }}
              </el-dropdown-item>
              <el-dropdown-item divided command="">全部产品</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        
        <!-- 解决方案下拉菜单 -->
        <el-dropdown trigger="hover" @command="(cmd: string) => $router.push(`/solutions/${cmd}`)">
          <span class="nav-link">
            解决方案<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="game">游戏加速</el-dropdown-item>
              <el-dropdown-item command="video">视频直播</el-dropdown-item>
              <el-dropdown-item command="edu">在线教育</el-dropdown-item>
              <el-dropdown-item command="ecommerce">电商平台</el-dropdown-item>
              <el-dropdown-item command="security">安全防护</el-dropdown-item>
              <el-dropdown-item divided command="">全部方案</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        
        <router-link to="/news" class="nav-link">新闻动态</router-link>
        <router-link to="/about" class="nav-link">关于我们</router-link>
        
        <!-- 帮助下拉菜单 -->
        <el-dropdown trigger="hover" @command="(cmd: string) => $router.push(`/${cmd}`)">
          <span class="nav-link">
            帮助支持<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="help">帮助中心</el-dropdown-item>
              <el-dropdown-item command="knowledge-base">知识库</el-dropdown-item>
              <el-dropdown-item command="downloads">下载中心</el-dropdown-item>
              <el-dropdown-item command="contact">联系我们</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </nav>
      
      <div class="header-actions">
        <template v-if="userStore.token">
          <el-dropdown trigger="click" @command="handleUserCommand">
            <div class="user-trigger">
              <el-avatar :size="32" class="user-avatar">{{ userInitial }}</el-avatar>
              <span class="username">{{ userStore.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="dashboard">控制台</el-dropdown-item>
                <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                <el-dropdown-item command="orders">我的订单</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <el-button text class="login-btn" @click="$router.push('/login')">
            <el-icon :size="16" style="margin-right: 4px;"><User /></el-icon>
            登录
          </el-button>
          <el-button type="primary" round size="default" @click="$router.push('/register')">
            免费注册
          </el-button>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import { ArrowDown, User } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

const scrolled = ref(false)
const productGroups = ref<any[]>([])
const siteName = ref('')

const fetchSiteName = async () => {
  try {
    const res = await request.get('/api/v1/settings/public')
    if (res?.data?.site_name) {
      siteName.value = res.data.site_name
    }
  } catch {
    // Use empty
  }
}

const userInitial = computed(() => {
  const username = userStore.username || 'U'
  return username.charAt(0).toUpperCase()
})

const handleScroll = () => {
  scrolled.value = window.scrollY > 50
}

const fetchProductGroups = async () => {
  try {
    const { data } = await request.get('/api/v1/product-groups')
    if (data?.data) {
      productGroups.value = data.data
    }
  } catch (error) {
    console.error('获取产品分组失败:', error)
  }
}

const handleUserCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  } else {
    router.push(`/user/${command}`)
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  fetchProductGroups()
  fetchSiteName()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="scss">
.site-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  
  &.header-scrolled {
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  }
  
  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    
    .logo-img {
      width: 36px;
      height: 36px;
      border-radius: 8px;
    }
    
    .logo-text {
      font-size: 18px;
      font-weight: 700;
      color: #1a2332;
    }
  }
  
  .nav-links {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .nav-link {
      padding: 8px 16px;
      color: #333;
      text-decoration: none;
      font-size: 14px;
      font-weight: 500;
      border-radius: 6px;
      transition: all 0.3s;
      cursor: pointer;
      
      &:hover {
        color: #409eff;
        background: rgba(64, 158, 255, 0.06);
      }
      
      &.router-link-active {
        color: #409eff;
      }
    }
  }
  
  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .login-btn {
      color: #333;
      
      &:hover {
        color: #409eff;
      }
    }
    
    .user-trigger {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 6px;
      transition: background 0.3s;
      
      &:hover {
        background: #f5f7fa;
      }
      
      .user-avatar {
        background: linear-gradient(135deg, #409eff, #66b1ff);
        color: #fff;
        font-weight: 600;
      }
      
      .username {
        font-size: 14px;
        color: #333;
      }
    }
  }
}

@media (max-width: 768px) {
  .site-header {
    .nav-links {
      display: none;
    }
  }
}
</style>
