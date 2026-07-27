<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <router-view />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import type { GlobalThemeOverrides } from 'naive-ui'

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
</style>
