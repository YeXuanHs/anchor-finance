<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <!-- 维护模式提示 -->
          <div v-if="isMaintenanceMode" class="maintenance-overlay">
            <div class="maintenance-content">
              <h2>系统维护中</h2>
              <p>{{ maintenanceMessage }}</p>
            </div>
          </div>
          <router-view v-else />
          <!-- AI导购浮窗 -->
          <AiShoppingWidget />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import type { GlobalThemeOverrides } from 'naive-ui'
import { useConfigStore } from '@/stores/config'
import AiShoppingWidget from '@/components/AiShoppingWidget.vue'

const configStore = useConfigStore()

const isMaintenanceMode = computed(() => configStore.config.maintenance_mode)
const maintenanceMessage = computed(() => '系统维护中，请稍后再访问')

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#1890ff',
    primaryColorHover: '#40a9ff',
    primaryColorPressed: '#096dd9',
    primaryColorSuppl: '#40a9ff',
    borderRadius: '12px',
    borderRadiusSmall: '8px'
  },
  Button: {
    borderRadiusMedium: '12px',
    borderRadiusLarge: '12px'
  },
  Card: {
    borderRadius: '12px'
  },
  Input: {
    borderRadius: '12px'
  },
  Tag: {
    borderRadius: '12px'
  }
}

onMounted(async () => {
  // 初始化时获取公开配置
  await configStore.fetchPublicConfig()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC',
    'Microsoft YaHei', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  min-height: 100vh;
}

.maintenance-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.maintenance-content {
  text-align: center;
  color: white;
  padding: 40px;
}

.maintenance-content h2 {
  font-size: 36px;
  margin-bottom: 20px;
}

.maintenance-content p {
  font-size: 18px;
  opacity: 0.9;
}
</style>
