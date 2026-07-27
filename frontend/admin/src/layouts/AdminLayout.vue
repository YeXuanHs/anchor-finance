<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <img src="/logo.png" alt="Logo" class="logo" />
        <span v-show="!isCollapsed" class="title">锚点财务</span>
      </div>
      
      <el-scrollbar class="sidebar-menu">
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapsed"
          :collapse-transition="false"
          router
          background-color="transparent"
          text-color="rgba(255, 255, 255, 0.85)"
          active-text-color="#ffffff"
        >
          <template v-for="route in menuRoutes" :key="route.path">
            <!-- 单级菜单 -->
            <el-menu-item
              v-if="!route.children || route.children.length === 1"
              :index="route.children ? route.children[0].path : route.path"
            >
              <el-icon><component :is="route.meta?.icon || 'Document'" /></el-icon>
              <template #title>{{ route.meta?.title }}</template>
            </el-menu-item>
            
            <!-- 多级菜单 -->
            <el-sub-menu v-else :index="route.path">
              <template #title>
                <el-icon><component :is="route.meta?.icon || 'Folder'" /></el-icon>
                <span>{{ route.meta?.title }}</span>
              </template>
              <el-menu-item
                v-for="child in route.children"
                :key="child.path"
                :index="`${route.path}/${child.path}`"
              >
                {{ child.meta?.title }}
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </el-scrollbar>
    </aside>
    
    <!-- 主内容区 -->
    <div class="main-area">
      <!-- 头部 -->
      <header class="header">
        <div class="header-left">
          <el-icon
            class="collapse-btn"
            :size="20"
            @click="isCollapsed = !isCollapsed"
          >
            <component :is="isCollapsed ? 'Expand' : 'Fold'" />
          </el-icon>
          
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute?.meta?.title">
              {{ currentRoute.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <el-tooltip content="刷新页面" placement="bottom">
            <el-icon class="header-icon" @click="refreshPage"><Refresh /></el-icon>
          </el-tooltip>
          
          <el-tooltip content="全屏" placement="bottom">
            <el-icon class="header-icon" @click="toggleFullscreen"><FullScreen /></el-icon>
          </el-tooltip>
          
          <el-badge :value="3" :max="99" class="notify-badge">
            <el-icon class="header-icon"><Bell /></el-icon>
          </el-badge>
          
          <el-dropdown trigger="click" @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" class="user-avatar">
                {{ userStore.username.charAt(0).toUpperCase() }}
              </el-avatar>
              <span class="username">{{ userStore.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人资料
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>系统设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      
      <!-- 内容区 -->
      <main class="content">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <keep-alive :include="cachedViews">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapsed = ref(false)
const cachedViews = ref<string[]>([])

const activeMenu = computed(() => route.path)
const currentRoute = computed(() => route)

const menuRoutes = computed(() => {
  const mainRoute = router.options.routes.find(r => r.path === '/')
  return mainRoute?.children?.filter(r => !r.meta?.hidden) || []
})

const refreshPage = () => {
  const { fullPath } = route
  router.replace('/redirect' + fullPath)
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      // 跳转个人资料
      break
    case 'settings':
      router.push('/system/general')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        userStore.logout()
        router.push('/login')
      } catch {
        // 取消
      }
      break
  }
}
</script>

<style scoped lang="scss">
.admin-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  overflow: hidden;
  
  &.collapsed {
    width: var(--sidebar-collapsed-width);
  }
  
  .sidebar-header {
    height: var(--header-height);
    display: flex;
    align-items: center;
    padding: 0 16px;
    gap: 12px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    
    .logo {
      width: 32px;
      height: 32px;
      flex-shrink: 0;
    }
    
    .title {
      font-size: 16px;
      font-weight: 600;
      color: #ffffff;
      white-space: nowrap;
    }
  }
  
  .sidebar-menu {
    flex: 1;
    padding: 8px 0;
    
    :deep(.el-menu) {
      border-right: none;
      
      .el-menu-item {
        height: 48px;
        line-height: 48px;
        margin: 2px 8px;
        border-radius: 8px;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1) !important;
        }
        
        &.is-active {
          background: rgba(255, 255, 255, 0.2) !important;
          font-weight: 500;
        }
      }
      
      .el-sub-menu {
        .el-sub-menu__title {
          height: 48px;
          line-height: 48px;
          margin: 2px 8px;
          border-radius: 8px;
          
          &:hover {
            background: rgba(255, 255, 255, 0.1) !important;
          }
        }
        
        .el-menu-item {
          min-width: auto;
          padding-left: 56px !important;
        }
      }
    }
  }
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-color);
}

.header {
  height: var(--header-height);
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: var(--shadow-sm);
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .collapse-btn {
      cursor: pointer;
      color: var(--text-secondary);
      transition: color 0.3s;
      
      &:hover {
        color: var(--primary-color);
      }
    }
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .header-icon {
      cursor: pointer;
      color: var(--text-secondary);
      transition: color 0.3s;
      
      &:hover {
        color: var(--primary-color);
      }
    }
    
    .notify-badge {
      display: flex;
      align-items: center;
    }
    
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 8px;
      transition: background 0.3s;
      
      &:hover {
        background: var(--bg-color);
      }
      
      .user-avatar {
        background: var(--primary-color);
        color: #fff;
        font-size: 14px;
      }
      
      .username {
        font-size: 14px;
        color: var(--text-primary);
      }
    }
  }
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
