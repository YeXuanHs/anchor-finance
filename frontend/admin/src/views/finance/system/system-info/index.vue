<template>
  <div class="system-info-page">
    <el-card shadow="never">
      <template #header><span>{{ $t('systemInfo.title') }}</span></template>

      <div class="system-info">
        <div class="info-header">
          <div class="logo"><el-icon :size="64" color="var(--el-color-primary)"><Monitor /></el-icon></div>
          <div class="title">
            <h1>{{ $t('systemInfo.systemName') }}</h1>
            <p class="version">{{ $t('systemInfo.version') }}: v{{ version }}</p>
            <p class="desc">{{ $t('systemInfo.description') }}</p>
          </div>
        </div>

        <el-divider />

        <el-descriptions :column="2" border :title="$t('systemInfo.systemInfoTitle')">
          <el-descriptions-item :label="$t('systemInfo.systemNameLabel')">{{ $t('systemInfo.systemNameValue') }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.currentVersion')">v{{ version }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.goVersion')">{{ systemInfo.go_version || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.dbVersion')">{{ systemInfo.db_version || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.os')">{{ systemInfo.os || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.serverTime')">{{ systemInfo.server_time || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.startedAt')">{{ systemInfo.started_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.uptime')">{{ systemInfo.uptime || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <el-descriptions :column="2" border :title="$t('systemInfo.licenseTitle')">
          <el-descriptions-item :label="$t('systemInfo.licenseType')"><el-tag type="success">{{ $t('systemInfo.openSource') }}</el-tag></el-descriptions-item>
          <el-descriptions-item :label="$t('systemInfo.licenseStatus')">{{ $t('systemInfo.permanent') }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <div class="links">
          <h3>{{ $t('systemInfo.relatedLinks') }}</h3>
          <el-space wrap>
            <el-button type="primary" plain @click="openLink('https://github.com/anchorfinance')">{{ $t('systemInfo.github') }}</el-button>
            <el-button type="primary" plain @click="openLink('https://docs.anchorfinance.dev')">{{ $t('systemInfo.docs') }}</el-button>
            <el-button type="primary" plain @click="checkUpdate">{{ $t('systemInfo.checkUpdate') }}</el-button>
          </el-space>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const version = ref('1.0.0')
const systemInfo = ref<any>({})

const fetchSystemInfo = async () => { try { const data = await request.get({ url: '/api/admin/system/info' }); systemInfo.value = data || {}; if (data?.version) version.value = data.version } catch {} }
const openLink = (url: string) => { window.open(url, '_blank') }

const checkUpdate = async () => {
  try { const data = await request.get({ url: '/api/admin/system/check-update' }); if (data?.has_update) ElMessage.info($t('systemInfo.updateAvailable', { version: data.latest_version })); else ElMessage.success($t('systemInfo.alreadyLatest')) } catch {}
}

onMounted(() => { fetchSystemInfo() })
</script>

<style scoped lang="scss">
.system-info-page { padding: 20px; }
.system-info { max-width: 800px; }
.info-header { display: flex; align-items: center; gap: 24px; margin-bottom: 20px; }
.title h1 { margin: 0 0 8px; font-size: 24px; }
.version { color: var(--el-text-color-secondary); margin: 0 0 4px; }
.desc { color: var(--el-text-color-secondary); margin: 0; font-size: 14px; }
.links { h3 { margin-bottom: 16px; } }
</style>
